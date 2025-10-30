import { memo } from "react";
import { ChatHistoryButton } from "./ChatHistoryButton";
import { ResetChatButton } from "./ResetChatButton";

interface ChatControlsProps {
  hasMessages: boolean;
  loading: boolean;
  onReset: () => void;
}

export const ChatControls = memo(
  ({ hasMessages, loading, onReset }: ChatControlsProps) => (
    <div className="flex items-center gap-2 sticky top-0 z-10 mb-4 justify-between px-4 sm:px-6 md:px-8 lg:px-12 xl:px-16 w-full animate-slide-in transition-all duration-300 ease-in-out">
      <ChatHistoryButton />
      {hasMessages && (
        <ResetChatButton
          isVisible={true}
          disabled={loading}
          onReset={onReset}
        />
      )}
    </div>
  )
);
