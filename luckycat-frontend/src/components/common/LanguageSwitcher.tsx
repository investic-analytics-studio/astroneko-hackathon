import { useState } from "react";
// import { useTranslation } from "react-i18next";
import { Globe, Check, ChevronDown } from "lucide-react";
import { useLanguage } from "@/hooks";
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover";
import { Language } from "@/hooks/common/useLanguage";

export const LanguageSwitcher = () => {
  // const { t } = useTranslation();
  const {
    currentLanguage,
    changeLanguage,
    getLanguageName,
    availableLanguages,
  } = useLanguage();
  const [isOpen, setIsOpen] = useState(false);

  const handleLanguageChange = (language: Language) => {
    changeLanguage(language);
    setIsOpen(false);
  };

  const getLanguageFlag = (language: Language): string => {
    const flags: Record<Language, string> = {
      en: "🇺🇸",
      th: "🇹🇭",
      jp: "🇯🇵",
    };
    return flags[language];
  };

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger asChild>
        <button
          className="flex items-center gap-1 sm:gap-2 px-2 sm:px-3 h-8 sm:h-9 md:h-10 lg:h-11 text-xs sm:text-sm font-medium text-white bg-transparent rounded-xl
        hover:bg-white/10 hover:text-white transition-all duration-300 cursor-pointer foucus:none ring-0 outline-none border-0"
        >
          <Globe className="hidden sm:block w-3 h-3 sm:w-4 sm:h-4" />
          {/* <span className="text-sm sm:text-base">
            {getLanguageFlag(currentLanguage)}
          </span> */}
          <span className="inline text-xs sm:text-sm">
            {getLanguageName(currentLanguage)}
          </span>
          <span>
            <ChevronDown className="w-3 h-3 sm:w-4 sm:h-4" />
          </span>
        </button>
      </PopoverTrigger>
      <PopoverContent
        className="w-auto p-1 bg-black/90 border-white/20 rounded-xl backdrop-blur-xl"
        align="end"
        sideOffset={8}
      >
        <div className="space-y-1">
          {/* <div className="px-2 sm:px-3 py-1.5 sm:py-2 text-xs font-semibold text-white/70 uppercase tracking-wider">
            {t("language.switch_language")}
          </div> */}
          {availableLanguages.map((language: Language) => (
            <button
              key={language}
              onClick={() => handleLanguageChange(language)}
              className={`w-full flex items-center gap-2 sm:gap-3 px-2 sm:px-3 py-1.5 sm:py-2 text-xs sm:text-sm rounded-md transition-all duration-200 focus:none ring-0 outline-none border-0 ${
                currentLanguage === language
                  ? "bg-white/20 text-white"
                  : "bg-transparent text-white/80 hover:bg-white/10 hover:text-white"
              }`}
            >
              <span className="text-base sm:text-lg">
                {getLanguageFlag(language)}
              </span>
              <span className="flex-1 text-left">
                {getLanguageName(language)}
              </span>
              {currentLanguage === language && (
                <div className="text-white animate-scale-in">
                  <Check className="w-3 h-3 sm:w-4 sm:h-4" />
                </div>
              )}
            </button>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  );
};
