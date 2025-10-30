import { MessageSquare, Trash2, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import type { ChatSession } from "@/apis/chat-history";
import { memo } from "react";
import { useTranslation } from "react-i18next";

interface ChatHistoryItemProps {
  session: ChatSession;
  onSelect: (sessionId: string, sessionName: string) => void;
  onDelete: (sessionId: string, e: React.MouseEvent) => void;
  isDeleting: boolean;
}

// Helper function to format timestamp
const formatTimestamp = (timestamp: string): string => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));
  const diffInDays = Math.floor(diffInHours / 24);
  const diffInWeeks = Math.floor(diffInDays / 7);
  const diffInMonths = Math.floor(diffInDays / 30);

  if (diffInHours < 1) return "Just now";
  if (diffInHours < 24)
    return `${diffInHours} hour${diffInHours > 1 ? "s" : ""} ago`;
  if (diffInDays < 7)
    return `${diffInDays} day${diffInDays > 1 ? "s" : ""} ago`;
  if (diffInWeeks < 4)
    return `${diffInWeeks} week${diffInWeeks > 1 ? "s" : ""} ago`;
  return `${diffInMonths} month${diffInMonths > 1 ? "s" : ""} ago`;
};

export const ChatHistoryItem = memo(
  ({ session, onSelect, onDelete, isDeleting }: ChatHistoryItemProps) => {
    const { t } = useTranslation();

    return (
      <div
        onClick={() => onSelect(session.session_id, session.history_name)}
        className="group relative p-4 rounded-lg border border-white/10 hover:border-[var(--brand-accent)]/30 hover:bg-[var(--brand-accent)]/5 transition-all duration-200 cursor-pointer"
      >
        <div className="flex items-start justify-between">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <MessageSquare className="h-4 w-4 text-[var(--brand-accent)] flex-shrink-0" />
              <h3 className="text-white font-press-start text-sm truncate">
                {session.history_name || t("chat.no_history_title_name")}
              </h3>
            </div>
            <p className="text-white/50 font-press-start text-xs">
              {formatTimestamp(session.updated_at)}
            </p>
          </div>

          <Button
            variant="ghost"
            size="icon"
            onClick={(e) => onDelete(session.session_id, e)}
            disabled={isDeleting}
            className="bg-gray-500/10 text-white/60 hover:text-[var(--brand-error)] hover:bg-[var(--brand-error)]/10 transition-all duration-200 h-6 w-6"
          >
            {isDeleting ? (
              <Loader2 className="h-3 w-3 animate-spin text-white" />
            ) : (
              <Trash2 className="h-3 w-3 text-white" />
            )}
            <span className="sr-only">Delete chat</span>
          </Button>
        </div>
      </div>
    );
  }
);

ChatHistoryItem.displayName = "ChatHistoryItem";
