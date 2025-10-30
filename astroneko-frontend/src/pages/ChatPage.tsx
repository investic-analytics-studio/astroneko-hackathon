import { useState, useEffect } from "react";
import { useChatController } from "@/hooks";
import {
  useChatHistoryStore,
  chatHistorySelectors,
} from "@/features/chat/history";
import { useChatData } from "@/features/chat/hooks/useChatData";
import { useChatStore } from "@/store/chatStore";
import { ChatLayout } from "@/features/chat/ui/ChatLayout";
import { ChatControls } from "@/features/chat/ui/ChatControls";
import { ChatMessages } from "@/features/chat/core/ChatMessages";
import { ChatInputArea } from "@/features/chat/ui/ChatInputArea";
import { TextOnlyLoadingSkeleton } from "@/features/chat/shared/LoadingSkeleton";

export default function ChatPage() {
  const {
    chatRef,
    input,
    setInput,
    messages,
    loading,
    showQuestions,
    categoryQuestions,
    toggleQuestions,
    handleSelectSuggestion,
    handleSendMessage,
    resetChat,
  } = useChatController();

  const currentSessionId = useChatHistoryStore(
    chatHistorySelectors.currentSessionId
  );
  const clearAllCategories = useChatStore((state) => state.clearAllCategories);

  // Clear current messages when switching to a different history session
  // This prevents messages from previous sessions from persisting
  useEffect(() => {
    if (currentSessionId) {
      clearAllCategories();
    }
  }, [currentSessionId, clearAllCategories]);

  // Check if chat history sidebar is open to hide suggestions
  const isHistorySidebarOpen = useChatHistoryStore(
    chatHistorySelectors.isHistorySidebarOpen
  );

  // Track if suggestions were hidden by sidebar opening
  const [suggestionsHiddenBySidebar, setSuggestionsHiddenBySidebar] =
    useState(false);

  // When sidebar opens, if suggestions were visible, hide them permanently
  useEffect(() => {
    if (isHistorySidebarOpen && showQuestions) {
      setSuggestionsHiddenBySidebar(true);
    }
  }, [isHistorySidebarOpen, showQuestions]);

  // Custom toggle function that resets the sidebar-hidden flag
  const handleToggleSuggestions = () => {
    setSuggestionsHiddenBySidebar(false);
    toggleQuestions();
  };

  const { messages: displayMessages, isLoading: isLoadingHistory } =
    useChatData(currentSessionId, messages);

  return (
    <ChatLayout
      chatRef={chatRef}
      chatContent={
        <>
          <ChatControls
            hasMessages={displayMessages.length > 0}
            loading={loading}
            onReset={() => void resetChat()}
          />

          <ChatMessages
            messages={displayMessages}
            isLoading={isLoadingHistory}
          />

          {loading && <TextOnlyLoadingSkeleton />}
        </>
      }
      chatInput={
        <ChatInputArea
          input={input}
          loading={loading}
          showSuggestions={
            showQuestions &&
            !isHistorySidebarOpen &&
            !suggestionsHiddenBySidebar
          }
          questions={categoryQuestions}
          onInputChange={setInput}
          onSendMessage={() => void handleSendMessage(input)}
          onToggleQuestions={handleToggleSuggestions}
          onSelectQuestion={handleSelectSuggestion}
        />
      }
    />
  );
}
