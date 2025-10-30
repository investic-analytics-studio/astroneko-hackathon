import { ArrowUp } from "lucide-react";
import { useRef } from "react";
import { useTranslation } from "react-i18next";

interface ChatInputProps {
  input: string;
  loading: boolean;
  onInputChange: (value: string) => void;
  onSendMessage: () => void;
  disabled?: boolean;
}

export const ChatInput = ({
  input,
  loading,
  onInputChange,
  onSendMessage,
  disabled = false,
}: ChatInputProps) => {
  const { t } = useTranslation();
  // Reference to the textarea element
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      onSendMessage();
    }
  };

  // Common classes for both input and textarea
  const commonClasses = `w-full px-4 2xl:px-12 py-2.5 rounded-2xl
    text-white placeholder-white/50
    border border-white/20 focus:border-white/25
    text-sm sm:text-base
    outline-none transition-all duration-300
    disabled:opacity-50 disabled:cursor-not-allowed
    scrollbar-hide`;

  return (
    <div className="sticky bottom-0 w-full px-4 bg-transparent transition-all duration-300 ease-in-out">
      <div className="max-w-4xl 2xl:max-w-[1700px] mx-auto relative">
        <textarea
          ref={textareaRef}
          rows={1}
          value={input}
          onChange={(e) => onInputChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={
            disabled ? t("chat.view_only_mode") : t("chat.ask_astroneko")
          }
          disabled={loading || disabled}
          className={`${commonClasses} bg-[#FFFFFF]/15 placeholder:flex placeholder:items-center resize-none min-h-[44px] sm:min-h-[48px] md:min-h-[48px] lg:min-h-[56px] 2xl:min-h-[64px] flex items-center rounded-[20px] sm:rounded-[22px] md:rounded-[24px] lg:rounded-[28px] 2xl:rounded-[32px] pr-[36px] sm:pr-[40px] md:pr-[50px] lg:pr-[60px] 2xl:pr-[80px] max-h-[150px] sm:max-h-[180px] md:max-h-[200px] lg:max-h-[250px] 2xl:max-h-[300px] overflow-y-auto`}
        />
        <button
          onClick={onSendMessage}
          disabled={loading || !input.trim() || disabled}
          className="absolute w-[32px] h-[32px] sm:w-[36px] sm:h-[36px] md:w-[40px] md:h-[40px] lg:w-[44px] lg:h-[44px] 2xl:w-[60px] 2xl:h-[60px] right-2 sm:right-2.5 md:right-3 lg:right-3.5 2xl:right-4 top-1/2 -translate-y-1/2 p-1.5 sm:p-2
                   bg-[#FFFFFF] hover:bg-[#FFFFFF] hover:border hover:border-white
                   rounded-full text-white hover:text-[#464749]
                   disabled:opacity-5 disabled:cursor-not-allowed
                   disabled:pointer-events-none
                   transition-all duration-300 ease-in-out
                   active:scale-95"
        >
          <ArrowUp
            className="group text-[#464749] disabled:opacity-40 w-4 h-4 sm:w-5 sm:h-5 md:w-5 md:h-5 lg:w-6 lg:h-6 2xl:w-8 2xl:h-8 mx-auto"
            strokeWidth={3}
          />
        </button>
      </div>
    </div>
  );
};
