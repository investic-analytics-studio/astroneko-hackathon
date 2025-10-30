import { Plus } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { useChatHistoryStore } from "@/features/chat/history";

export function NewChatButton() {
  const { t } = useTranslation();
  const startNewChat = useChatHistoryStore((state) => state.startNewChat);

  return (
    <Button
      onClick={startNewChat}
      className="bg-[var(--brand-primary)] hover:bg-[var(--brand-primary-hover)] text-white border-0 px-4 py-2 rounded-lg font-press-start text-xs transition-all duration-200 flex items-center gap-2 shadow-lg hover:shadow-xl"
    >
      <Plus className="h-4 w-4" />
      <span className="hidden sm:inline">{t("chat.new_chat")}</span>
      <span className="sm:hidden">New</span>
    </Button>
  );
}
