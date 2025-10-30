import { memo } from "react";
import { useAuth } from "@/hooks";
import { LanguageSwitcher } from "./LanguageSwitcher";
import { LogoButton } from "./LogoButton";
import { UserSection } from "./UserSection";
import { SocialLinks } from "./SocialLinks";
import { ChatNavigation } from "./ChatNavigation";
import { useLocation } from "@tanstack/react-router";
import { useLanguage } from "@/hooks";
import { GenesisButton } from "./GenesisButton";

interface HeaderProps {
  onScrollToHome?: () => void;
}

export const Header = memo(({ onScrollToHome }: HeaderProps) => {
  const { authUser, logout } = useAuth();
  const { currentLanguage } = useLanguage();
  const pathWithoutLang = useLocation().pathname.replace(
    `/${currentLanguage}`,
    ""
  );

  return (
    <div className="fixed top-0 left-0 right-0 py-0 md:px-4 lg:px-6 xl:px-8 z-50">
      <div className="md:hidden flex items-center justify-between bg-black/50 px-3 py-3">
        <div className="text-whites font-bold flex items-center gap-2">
          <span className="bg-gradient-to-r from-red-500 to-pink-600 text-[12px] md:text-[16px] text-white rounded-full px-2 py-1 text-xs shadow-md">
            New
          </span>
          <span className="text-[12px] md:text-[16px]">ASTRO GENESIS BOX</span>
        </div>
        <div>
          <GenesisButton />
        </div>
      </div>
      <div className="bg-[#e78562]/20 md:bg-transparent flex items-center justify-between px-3 lg:px-6 xl:px-8 2xl:px-12 py-3">
        {/* Left section: Logo */}
        <div>
          <LogoButton onScrollToHome={onScrollToHome} />
        </div>

        {/* Center section: Chat Navigation */}
        {pathWithoutLang.startsWith("/chat") && (
          <ChatNavigation pathWithoutLang={pathWithoutLang} />
        )}

        {/* Right section: Language, User, Social */}
        <div className="flex-shrink-0 font-press-start flex items-center justify-end gap-1.5 sm:gap-3 md:gap-3">
          <div className="hidden md:block">
            <GenesisButton />
          </div>
          <LanguageSwitcher />
          <UserSection
            authUser={authUser}
            isLoading={false}
            onLogout={logout}
          />
          <SocialLinks />
        </div>
      </div>
    </div>
  );
});
