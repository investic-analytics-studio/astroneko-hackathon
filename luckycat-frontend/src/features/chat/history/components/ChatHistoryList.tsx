import { Loader2, MessageSquare } from "lucide-react";
import { useTranslation } from "react-i18next";
import { ChatHistoryItem } from "./ChatHistoryItem";
import { ChatHistoryEmpty } from "./ChatHistoryEmpty";
import type { ChatSession } from "@/apis/chat-history";
import { memo } from "react";

interface ChatHistoryListProps {
  sessions: ChatSession[];
  isLoading: boolean;
  isError: boolean;
  searchQuery: string;
  onSelectSession: (sessionId: string, sessionName: string) => void;
  onDeleteSession: (sessionId: string, e: React.MouseEvent) => void;
  isDeletingSession: boolean;
}

export const ChatHistoryList = memo(
  ({
    sessions,
    isLoading,
    isError,
    searchQuery,
    onSelectSession,
    onDeleteSession,
    isDeletingSession,
  }: ChatHistoryListProps) => {
    const { t } = useTranslation();

    if (isLoading) {
      return (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-6 w-6 text-white/60 animate-spin" />
          <span className="ml-2 text-white/60 font-press-start text-sm">
            {t("chat.loading_conversations")}
          </span>
        </div>
      );
    }

    if (isError) {
      return (
        <div className="text-center py-8">
          <MessageSquare className="h-12 w-12 text-red-400 mx-auto mb-3" />
          <p className="text-red-400 font-press-start text-sm">
            {t("chat.failed_to_load")}
          </p>
          <p className="text-white/50 font-press-start text-xs mt-1">
            Please try again later
          </p>
        </div>
      );
    }

    if (sessions.length === 0) {
      return (
        <ChatHistoryEmpty
          message={
            searchQuery
              ? t("chat.no_conversations_found")
              : t("chat.no_conversations")
          }
          isSearching={!!searchQuery}
        />
      );
    }

    return (
      <div className="space-y-3">
        {sessions.map((session) => (
          <ChatHistoryItem
            key={session.session_id}
            session={session}
            onSelect={onSelectSession}
            onDelete={onDeleteSession}
            isDeleting={isDeletingSession}
          />
        ))}
      </div>
    );
  }
);

ChatHistoryList.displayName = "ChatHistoryList";
