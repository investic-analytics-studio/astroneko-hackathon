import { useCallback } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useLanguage } from "@/hooks";
import { Cat, DollarSign, Heart } from "lucide-react";

// Chat navigation button component
interface ChatNavButtonProps {
  isActive: boolean;
  onClick: () => void;
  onKeyDown: (event: React.KeyboardEvent) => void;
  children: React.ReactNode;
  ariaLabel: string;
}

const ChatNavButton = ({
  isActive,
  onClick,
  onKeyDown,
  children,
  ariaLabel,
}: ChatNavButtonProps) => {
  return (
    <button
      type="button"
      className={`p-[6px] sm:p-[8px] lg:p-[10px] xl:p-[12px] flex items-center justify-center border border-white/50 rounded-md hover:bg-white/20 cursor-pointer transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-white/50 text-white bg-white/10 ${
        isActive ? "bg-yellow-400/20 text-yellow-400 border-yellow-400" : ""
      }`}
      onClick={onClick}
      onKeyDown={onKeyDown}
      aria-label={ariaLabel}
      tabIndex={0}
    >
      <div>{children}</div>
    </button>
  );
};

interface ChatNavigationProps {
  pathWithoutLang: string;
}

export const ChatNavigation = ({ pathWithoutLang }: ChatNavigationProps) => {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { currentLanguage } = useLanguage();

  const handleGeneralChatClick = useCallback(() => {
    navigate({
      to: "/$lng/chat/$category",
      params: {
        lng: currentLanguage,
        category: "general",
      },
    });
  }, [navigate, currentLanguage]);

  const handleCryptoChatClick = useCallback(() => {
    navigate({
      to: "/$lng/chat/$category",
      params: {
        lng: currentLanguage,
        category: "crypto",
      },
    });
  }, [navigate, currentLanguage]);

  const handleLoverChatClick = useCallback(() => {
    navigate({
      to: "/$lng/chat/$category",
      params: {
        lng: currentLanguage,
        category: "lover",
      },
    });
  }, [navigate, currentLanguage]);

  const handleChatKeyDown = useCallback(
    (event: React.KeyboardEvent, handler: () => void) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        handler();
      }
    },
    []
  );

  // Don't render if on home or category page
  if (pathWithoutLang === "/" || pathWithoutLang === "/category") {
    return null;
  }

  return (
    <nav
      className="absolute left-1/2 transform -translate-x-1/2 flex items-center justify-center gap-1 sm:gap-2 lg:gap-3 xl:gap-4"
      role="navigation"
      aria-label="Chat categories"
    >
      <ChatNavButton
        isActive={pathWithoutLang === "/chat/general"}
        onClick={handleGeneralChatClick}
        onKeyDown={(e) => handleChatKeyDown(e, handleGeneralChatClick)}
        ariaLabel={t("common.general_chat")}
      >
        <Cat className="w-4 h-4 md:w-5 md:h-5" />
      </ChatNavButton>
      <ChatNavButton
        isActive={pathWithoutLang === "/chat/crypto"}
        onClick={handleCryptoChatClick}
        onKeyDown={(e) => handleChatKeyDown(e, handleCryptoChatClick)}
        ariaLabel={t("common.crypto_chat")}
      >
        <DollarSign className="w-4 h-4 md:w-5 md:h-5" />
      </ChatNavButton>
      <ChatNavButton
        isActive={pathWithoutLang === "/chat/lover"}
        onClick={handleLoverChatClick}
        onKeyDown={(e) => handleChatKeyDown(e, handleLoverChatClick)}
        ariaLabel={t("common.lover_chat")}
      >
        <Heart className="w-4 h-4 md:w-5 md:h-5" />
      </ChatNavButton>
    </nav>
  );
};
