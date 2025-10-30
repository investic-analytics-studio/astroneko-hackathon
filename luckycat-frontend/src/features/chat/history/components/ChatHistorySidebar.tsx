import { X, Plus } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { toast } from "sonner";
import { useAuth } from "@/hooks";
import { useChatStore } from "@/store/chatStore";
import {
  useChatSessions,
  useDeleteChatSession,
  useChatHistoryStore,
  chatHistorySelectors,
} from "@/features/chat/history";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { ChatHistorySearch } from "./ChatHistorySearch";
import { ChatHistoryList } from "./ChatHistoryList";
import { DeleteSessionModal } from "./DeleteSessionModal";

export function ChatHistorySidebar() {
  const { isAuthenticated, authUser } = useAuth();
  const { t } = useTranslation();

  // Store state
  const isOpen = useChatHistoryStore(chatHistorySelectors.isHistorySidebarOpen);
  const setIsOpen = useChatHistoryStore((state) => state.setHistorySidebarOpen);
  const searchQuery = useChatHistoryStore(chatHistorySelectors.searchQuery);
  const setSearchQuery = useChatHistoryStore((state) => state.setSearchQuery);
  const viewHistory = useChatHistoryStore((state) => state.viewHistory);
  const startNewChat = useChatHistoryStore((state) => state.startNewChat);

  // API hooks - only fetch sessions when authenticated
  const {
    data: sessionsData,
    isLoading: isLoadingSessions,
    error: sessionsError,
  } = useChatSessions(
    {
      sort_by: "updated_at",
      sort_order: "desc",
    },
    {
      enabled: isAuthenticated && !!authUser, // Only fetch when authenticated
    }
  );

  const deleteSessionMutation = useDeleteChatSession();

  // Delete modal state
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [sessionToDelete, setSessionToDelete] = useState<{
    id: string;
    name: string;
  } | null>(null);

  // Sync API data to store for client-side filtering
  useEffect(() => {
    if (sessionsData?.sessions) {
      useChatHistoryStore.getState().setSessions(sessionsData.sessions);
    }
  }, [sessionsData]);

  // Event handlers
  const handleChatSelect = (sessionId: string, sessionName: string) => {
    viewHistory(sessionId, sessionName);
    setIsOpen(false);
  };

  const handleDeleteChat = (sessionId: string, e: React.MouseEvent) => {
    e.stopPropagation();

    // Find session name
    const session = sessions.find((s) => s.session_id === sessionId);
    const sessionName = session?.history_name || "Unknown";

    setSessionToDelete({ id: sessionId, name: sessionName });
    setDeleteModalOpen(true);
  };

  const handleConfirmDelete = async () => {
    if (!sessionToDelete) return;

    const currentSessionId = useChatHistoryStore.getState().currentSessionId;
    const isCurrentSession = currentSessionId === sessionToDelete.id;

    await deleteSessionMutation.mutateAsync(sessionToDelete.id);

    // If deleting current session, clear messages and start new chat
    if (isCurrentSession) {
      const { startNewChat } = useChatHistoryStore.getState();
      const { clearAllCategories } = useChatStore.getState();

      // Clear all chat messages
      clearAllCategories();

      // Start new chat
      startNewChat();

      // Show toast for current session deletion
      toast.success(t("chat.session_deleted_and_refreshed"), {
        className: "custom-success-toast",
        descriptionClassName: "text-[#A1A1AA]",
        actionButtonStyle: { color: "#000000" },
      });
    } else {
      // Show toast for other session deletion
      toast.success(t("chat.session_deleted"), {
        className: "custom-success-toast",
        descriptionClassName: "text-[#A1A1AA]",
        actionButtonStyle: { color: "#000000" },
      });
    }

    setDeleteModalOpen(false);
    setSessionToDelete(null);
  };

  const handleNewChat = () => {
    // Clear all chat messages
    const { clearAllCategories } = useChatStore.getState();
    clearAllCategories();

    // Start new chat
    startNewChat();
    setIsOpen(false);
  };

  // Get sessions from store (client-side filtering)
  const sessions = useChatHistoryStore(chatHistorySelectors.sessions);

  // Memoize filtered sessions to prevent infinite re-renders
  const filteredSessions = useMemo(() => {
    if (!searchQuery.trim()) return sessions;

    // Escape special regex characters to prevent errors
    const escapeRegex = (str: string) => {
      return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    };

    const normalizedQuery = escapeRegex(searchQuery.toLowerCase().trim());
    return sessions.filter((session) =>
      session.history_name.toLowerCase().includes(normalizedQuery)
    );
  }, [sessions, searchQuery]);

  return (
    <>
      <Sheet open={isOpen} onOpenChange={setIsOpen}>
        <SheetContent
          side="left"
          className="bg-black/40 backdrop-blur-xl border-r border-white/20 w-full sm:w-80 md:w-96 p-0 [&>button]:hidden flex flex-col h-full"
        >
          <SheetHeader className="p-4 sm:p-6 pb-4 border-b border-white/10">
            <div className="flex items-center justify-between mb-4">
              <SheetTitle className="text-white font-press-start text-xl">
                Chat History
              </SheetTitle>

              <div className="flex items-center gap-2">
                {/* New Chat Button */}
                <Button
                  onClick={handleNewChat}
                  className="bg-white/20 hover:bg-white/25 text-white border border-white/20 hover:border-white/30 px-3 py-2 rounded-lg font-press-start text-xs transition-all duration-200 flex items-center gap-2 h-8 shadow-lg backdrop-blur-sm"
                >
                  <Plus className="h-3 w-3" />
                  <span>{t("chat.new_chat")}</span>
                </Button>

                {/* Custom Close Button */}
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setIsOpen(false)}
                  className="bg-white/20 hover:bg-white/25 text-white border border-white/20 hover:border-white/30 hover:text-white rounded-lg transition-all duration-200 h-8 w-8 shadow-lg backdrop-blur-sm"
                >
                  <X className="h-4 w-4" />
                  <span className="sr-only">Close</span>
                </Button>
              </div>
            </div>

            {/* Search Input - Real-time client-side filtering */}
            {isAuthenticated && (
              <ChatHistorySearch
                value={searchQuery}
                onChange={setSearchQuery}
              />
            )}
          </SheetHeader>

          <div className="flex flex-col flex-1 min-h-0">
            {/* Chat History List - Only for authenticated users */}
            {isAuthenticated && (
              <div className="flex-1 p-3 sm:p-4 overflow-y-auto min-h-0">
                <div className="space-y-3">
                  <div className="text-white/60 text-sm font-press-start mb-4">
                    {searchQuery
                      ? `${t("chat.search_results")} (${
                          filteredSessions.length
                        })`
                      : t("chat.recent_conversations")}
                  </div>

                  <ChatHistoryList
                    sessions={filteredSessions}
                    isLoading={isLoadingSessions}
                    isError={!!sessionsError}
                    searchQuery={searchQuery}
                    onSelectSession={handleChatSelect}
                    onDeleteSession={handleDeleteChat}
                    isDeletingSession={deleteSessionMutation.isPending}
                  />
                </div>
              </div>
            )}

            {/* Footer */}
            <div className="p-3 sm:p-4 border-t border-white/10">
              <div className="text-center text-white/50 text-xs font-press-start">
                {isAuthenticated
                  ? `${filteredSessions.length} of ${sessions.length} conversations`
                  : t("chat.sign_in_to_save")}
              </div>
            </div>
          </div>
        </SheetContent>
      </Sheet>

      {/* Delete Confirmation Modal */}
      <DeleteSessionModal
        isOpen={deleteModalOpen}
        onOpenChange={setDeleteModalOpen}
        sessionName={sessionToDelete?.name || ""}
        onConfirmDelete={handleConfirmDelete}
        isDeleting={deleteSessionMutation.isPending}
      />
    </>
  );
}
