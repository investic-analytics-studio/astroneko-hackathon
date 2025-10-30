import { useCallback, useEffect, memo } from "react";
import { Check, Copy } from "lucide-react";
import { useTranslation } from "react-i18next";

import { FormatText, FortuneCard } from "./FortuneCard";
import { useClipboard } from "@/hooks";
import { FortuneMessage } from "@/hooks/fortune/useFortune";
import { logger } from "@/lib/logger";

interface ChatMessageProps {
  message: FortuneMessage;
}

// Custom comparison to prevent re-renders for same message ID
const arePropsEqual = (
  prevProps: ChatMessageProps,
  nextProps: ChatMessageProps
) => {
  return prevProps.message.id === nextProps.message.id;
};

export const ChatMessage = memo(({ message }: ChatMessageProps) => {
  const { copied, error, copy } = useClipboard();
  const { t } = useTranslation();

  useEffect(() => {
    if (error) {
      logger.error("Failed to copy message text", error);
    }
  }, [error]);

  const handleCopy = useCallback(
    (value: string) => {
      void copy(value);
    },
    [copy]
  );
  return (
    <div
      key={message.id}
      className={`flex mx-4 md:mx-20 opacity-0 translate-y-5 animate-fade-in ${
        message.role === "user" ? "justify-end text-left" : "justify-start"
      }`}
    >
      <div
        className={`relative ${
          message.role === "user"
            ? "max-w-[300px] xl:max-w-[1000px]"
            : message.card
            ? ""
            : ""
        } px-0 py-2 space-y-20 shadow-none ${
          message.role === "user"
            ? "bg-[#BD042D] text-white text-[14px] sm:text-[15px] md:text-[16px] lg:text-[17px] xl:text-[16px] 2xl:text-[18px] font-normal rounded-4xl rounded-tr-none px-4 sm:px-5 md:px-6 lg:px-7 xl:px-6 2xl:px-8"
            : "bg-transparent border-none text-white text-[15px] sm:text-[16px] md:text-[17px] lg:text-[18px] xl:text-[17px] 2xl:text-[19px] font-normal xl:font-medium rounded-[14px] 2xl:rounded-[30px]"
        }`}
      >
        {message.role === "ai" ? (
          <div className="space-y-3 tracking-wide leading-relaxed text-start px-5">
            {message.card && message.meaning ? (
              <FortuneCard
                card={message.card}
                meaning={message.meaning}
                text={message.message}
              />
            ) : (
              <div className="relative text-[15px] sm:text-[16px] md:text-[17px] lg:text-[18px] xl:text-[17px] 2xl:text-[19px] 2xl:px-3">
                {FormatText(message.message)}
              </div>
            )}
            {message.role === "ai" && (
              <button
                type="button"
                onClick={() => handleCopy(message.message)}
                className="text-[11px] sm:text-[12px] md:text-[13px] lg:text-[14px] xl:text-[13px] 2xl:text-[15px] p-1.5 sm:p-2 rounded-md border border-white/20 bg-white/5 hover:bg-white/10 hover:border-white/50 text-white transition-colors focus:outline-none"
                title={copied ? t("chat.copied") : t("chat.copy_raw_text")}
                aria-label={t("chat.copy_message")}
              >
                {copied ? (
                  <Check className="w-4 h-4 text-green-500" />
                ) : (
                  <Copy className="w-4 h-4" />
                )}
              </button>
            )}
          </div>
        ) : (
          <div className="relative">
            {message.message}
            <div className="absolute -inset-1 rounded-lg blur-sm animate-pulse-glow" />
          </div>
        )}
      </div>
    </div>
  );
}, arePropsEqual);
