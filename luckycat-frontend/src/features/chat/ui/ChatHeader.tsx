import { ChatHistoryButton } from "./ChatHistoryButton";
import { NewChatButton } from "./NewChatButton";
import { ResetChatButton } from "./ResetChatButton";
import {
  useChatHistoryStore,
  chatHistorySelectors,
} from "@/features/chat/history";

interface ChatHeaderProps {
  onReset: () => void;
  resetDisabled: boolean;
  showReset: boolean;
}

export function ChatHeader({
  onReset,
  resetDisabled,
  showReset,
}: ChatHeaderProps) {
  const chatMode = useChatHistoryStore(chatHistorySelectors.chatMode);
  const isViewingHistory = chatMode === "history";

  return (
    <div className="flex gap-2 sticky top-0 z-10 mb-4 justify-between px-4 sm:px-6 md:px-8 lg:px-12 xl:px-16 w-full">
      <div className="flex gap-2">
        <ChatHistoryButton />
        {showReset && !isViewingHistory && (
          <ResetChatButton
            isVisible={showReset}
            disabled={resetDisabled}
            onReset={onReset}
          />
        )}
      </div>

      {chatMode !== "new" && (
        <div className="flex gap-2">
          <NewChatButton />
        </div>
      )}
    </div>
  );
}
