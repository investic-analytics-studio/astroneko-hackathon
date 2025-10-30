import { useEffect, useRef } from "react";
import { Loader2, MessageSquare } from "lucide-react";
import { ChatMessage } from "@/features/chat/core";
import {
  useChatHistory,
  useChatHistoryStore,
  useRefetchChatHistory,
} from "@/features/chat/history";
import type { FortuneMessage } from "@/hooks/fortune/useFortune";
import { convertApiMessagesToFortuneMessages } from "../utils/chat-history-utils";

interface ChatHistoryViewerProps {
  className?: string;
}

export function ChatHistoryViewer({ className = "" }: ChatHistoryViewerProps) {
  const chatRef = useRef<HTMLDivElement>(null);

  const currentSessionId = useChatHistoryStore(
    (state) => state.currentSessionId
  );
  const refetchChatHistory = useRefetchChatHistory();

  // API hook
  const {
    data: historyData,
    isLoading,
    error,
  } = useChatHistory(currentSessionId || "", !!currentSessionId);

  // Refetch when session changes to ensure fresh data
  useEffect(() => {
    if (currentSessionId) {
      refetchChatHistory(currentSessionId);
    }
  }, [currentSessionId, refetchChatHistory]);

  // Auto-scroll to bottom when new messages load
  useEffect(() => {
    if (chatRef.current && historyData?.messages) {
      chatRef.current.scrollTop = chatRef.current.scrollHeight;
    }
  }, [historyData?.messages]);

  // Convert API messages to FortuneMessage format
  const messages = historyData?.messages
    ? convertApiMessagesToFortuneMessages(historyData.messages)
    : [];

  if (!currentSessionId) {
    return (
      <div className={`flex-1 flex items-center justify-center ${className}`}>
        <div className="text-center">
          <MessageSquare className="h-16 w-16 text-white/30 mx-auto mb-4" />
          <h3 className="text-white/60 font-press-start text-lg mb-2">
            No Session Selected
          </h3>
          <p className="text-white/40 font-press-start text-sm">
            Select a conversation from the history to view messages
          </p>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className={`flex-1 flex items-center justify-center ${className}`}>
        <div className="flex items-center gap-3">
          <Loader2 className="h-6 w-6 text-white/60 animate-spin" />
          <span className="text-white/60 font-press-start text-sm">
            Loading conversation...
          </span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`flex-1 flex items-center justify-center ${className}`}>
        <div className="text-center">
          <MessageSquare className="h-16 w-16 text-red-400 mx-auto mb-4" />
          <h3 className="text-red-400 font-press-start text-lg mb-2">
            Failed to Load
          </h3>
          <p className="text-white/60 font-press-start text-sm">
            Could not load conversation messages
          </p>
        </div>
      </div>
    );
  }

  if (!messages.length) {
    return (
      <div className={`flex-1 flex items-center justify-center ${className}`}>
        <div className="text-center">
          <MessageSquare className="h-16 w-16 text-white/30 mx-auto mb-4" />
          <h3 className="text-white/60 font-press-start text-lg mb-2">
            No Messages
          </h3>
          <p className="text-white/40 font-press-start text-sm">
            This conversation appears to be empty
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={`flex flex-col h-full ${className}`}>
      {/* Messages Container */}
      <div ref={chatRef} className="flex-1 overflow-y-auto px-4 py-4 space-y-4">
        {messages.map((message: FortuneMessage, index: number) => (
          <div
            key={message.id}
            className="opacity-0 translate-y-5 animate-fade-in"
            style={{
              animationDelay: `${index * 0.1}s`,
            }}
          >
            <ChatMessage message={message} />
          </div>
        ))}
      </div>
    </div>
  );
}
