import { Trash } from "lucide-react";
import { useTranslation } from "react-i18next";

interface ResetChatButtonProps {
  isVisible: boolean;
  disabled: boolean;
  onReset: () => void;
}

export const ResetChatButton = ({
  isVisible,
  disabled,
  onReset,
}: ResetChatButtonProps) => {
  const { t } = useTranslation();

  return (
    <>
      {isVisible && (
        <button
          onClick={onReset}
          disabled={disabled}
          className={`flex items-center justify-center gap-1 rounded-md border text-sm transition-all duration-200 shadow-lg bg-gray-400/90 text-white border-white/20 hover:border-white/30 hover:bg-gray-400/70 focus:outline-none h-8 px-2 sm:h-9 sm:px-3 md:h-9 lg:h-10 ${
            disabled
              ? "cursor-not-allowed bg-white/10 text-white/50 border-white/10"
              : ""
          }`}
        >
          {disabled ? (
            <span className="font-press-start text-[9px] md:text-[10px] lg:text-[11px]">
              {t("chat.processing")}
            </span>
          ) : (
            <>
              <Trash className="h-3 w-3 md:h-4 md:w-4 lg:h-4 lg:w-4" />
              <span className="font-press-start text-[9px] md:text-[10px] lg:text-[11px]">
                {t("chat.clear_messages")}
              </span>
            </>
          )}
          <span className="sr-only">{t("chat.clear_messages")}</span>
        </button>
      )}
    </>
  );
};

ResetChatButton.displayName = "ResetChatButton";
