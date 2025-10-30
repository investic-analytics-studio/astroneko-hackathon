import { useTranslation } from "react-i18next";

export function HeroSection() {
  const { t } = useTranslation();

  return (
    <div className="opacity-0 translate-y-5 animate-fade-in">
      <div className="flex items-center justify-center mb-0 sm:mb-2 md:mb-6">
        <img
          src="/logo/astro-logo.webp"
          alt="logo"
          className="w-[280px] sm:w-[320px] md:w-[400px] lg:w-[500px] h-auto"
        />
      </div>
      <div
        className="text-xs sm:text-base md:text-lg lg:text-xl mx-auto text-white font-semibold font-press-start leading-relaxed py-2 space-y-1 md:space-y-2 opacity-0 translate-y-5 animate-fade-in"
        style={{
          textShadow:
            "4px 4px 8px rgba(0, 0, 0, 0.4), 2px 2px 4px rgba(0, 0, 0, 0.4)",
          animationDelay: "0.2s",
        }}
      >
        <p>{t("hero.title")}</p>
        <p>{t("hero.subtitle")}</p>
      </div>
    </div>
  );
}
