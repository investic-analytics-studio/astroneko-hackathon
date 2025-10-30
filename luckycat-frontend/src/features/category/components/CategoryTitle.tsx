import { useTranslation } from "react-i18next";

export function CategoryTitle() {
  const { t } = useTranslation();

  return (
    <h1
      className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl xl:text-6xl font-bold text-white font-press-start mb-6 sm:mb-8 md:mb-10 lg:mb-12 opacity-0 translate-y-5 animate-fade-in text-center px-4"
      style={{
        textShadow: "2px 2px 4px rgba(0, 0, 0, 0.3)",
      }}
    >
      {t("category.title")}
    </h1>
  );
}
