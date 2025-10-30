import { History } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/hooks";
import { useChatHistoryStore } from "@/features/chat/history";

export function ChatHistoryButton() {
  const { t } = useTranslation();
  const { isAuthenticated } = useAuth();
  const setHistorySidebarOpen = useChatHistoryStore(
    (state) => state.setHistorySidebarOpen
  );
  const isHistorySidebarOpen = useChatHistoryStore(
    (state) => state.isHistorySidebarOpen
  );

  const toggleChatHistory = () => setHistorySidebarOpen(!isHistorySidebarOpen);

  // Don't render if user is not authenticated
  if (!isAuthenticated) {
    return null;
  }

  return (
    <button
      onClick={toggleChatHistory}
      className="flex items-center justify-center gap-1 rounded-md border text-sm transition-all duration-200 shadow-lg bg-gray-400/90 text-white border-white/20 hover:border-white/30 hover:bg-gray-400/70 h-8 px-2 sm:h-9 sm:px-3 md:h-9 lg:h-10 focus:outline-none"
    >
      <History className="h-3 w-3 md:h-4 md:w-4 lg:h-4 lg:w-4" />
      <span className="font-press-start text-[9px] md:text-[10px] lg:text-[11px]">
        {t("chat.chat_history")}
      </span>
      <span className="sr-only">{t("chat.open_chat_history")}</span>
    </button>
  );
}
