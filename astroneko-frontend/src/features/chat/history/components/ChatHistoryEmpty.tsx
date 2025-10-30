import { MessageSquare } from "lucide-react";
import { memo } from "react";

interface ChatHistoryEmptyProps {
  message: string;
  isSearching?: boolean;
}

export const ChatHistoryEmpty = memo(
  ({ message, isSearching = false }: ChatHistoryEmptyProps) => {
    return (
      <div className="text-center py-8">
        <MessageSquare className="h-12 w-12 text-white/30 mx-auto mb-3" />
        <p className="text-white/50 font-press-start text-sm">{message}</p>
        {isSearching && (
          <p className="text-white/30 font-press-start text-xs mt-2">
            Try different keywords
          </p>
        )}
      </div>
    );
  }
);

ChatHistoryEmpty.displayName = "ChatHistoryEmpty";
