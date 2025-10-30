import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "@tanstack/react-router";

export type Language = "en" | "th" | "jp";

export const useLanguage = () => {
  const { i18n } = useTranslation();
  const navigate = useNavigate();
  const params = useParams({ from: "/$lng" }) as { lng?: string };
  const [currentLanguage, setCurrentLanguage] = useState<Language>("en");

  useEffect(() => {
    // Get language from URL params first, then localStorage, then default to 'en'
    const urlLanguage = params.lng as Language;
    const savedLanguage = localStorage.getItem("language") as Language;

    if (urlLanguage && ["en", "th", "jp"].includes(urlLanguage)) {
      setCurrentLanguage(urlLanguage);
      i18n.changeLanguage(urlLanguage);
      localStorage.setItem("language", urlLanguage);
    } else if (savedLanguage && ["en", "th", "jp"].includes(savedLanguage)) {
      setCurrentLanguage(savedLanguage);
      i18n.changeLanguage(savedLanguage);
    }
  }, [i18n, params.lng]);

  const changeLanguage = (language: Language) => {
    setCurrentLanguage(language);
    i18n.changeLanguage(language);
    localStorage.setItem("language", language);

    // Navigate to the same route but with new language
    const currentPath = window.location.pathname;
    const pathWithoutLang = currentPath.replace(/^\/[a-z]{2}/, "");
    const newPath = `/${language}${pathWithoutLang}`;

    navigate({ to: newPath });
  };

  const getLanguageName = (language: Language): string => {
    const names = {
      en: "EN",
      th: "TH",
      jp: "JP",
    };
    return names[language];
  };

  return {
    currentLanguage,
    changeLanguage,
    getLanguageName,
    availableLanguages: ["en", "th", "jp"] as Language[],
  };
};
