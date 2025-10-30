import { ArrowRight } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";

/**
 * @deprecated This component is no longer needed as users can directly type to continue conversations.
 * Kept for backward compatibility but does nothing.
 */
export function ContinueChatButton() {
  const { t } = useTranslation();

  return (
    <Button
      onClick={() => {
        // No-op: Users can directly type to continue conversations
      }}
      className="bg-[var(--brand-accent)] hover:bg-[var(--brand-accent)]/80 text-white border-0 px-4 py-2 rounded-lg font-press-start text-xs transition-all duration-200 flex items-center gap-2 shadow-lg hover:shadow-xl"
    >
      <span>{t("chat.continue_conversation")}</span>
      <ArrowRight className="h-4 w-4" />
    </Button>
  );
}
