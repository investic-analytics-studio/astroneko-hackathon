import { Cat, Diamond, DollarSign, Heart } from "lucide-react";
import { TFunction } from "i18next";

// Initial questions for each category
export const getInitialQuestions = (t: TFunction) => ({
  general:
    (t("chat.initial_questions.general", {
      returnObjects: true,
    }) as string[]) || [],
  crypto:
    (t("chat.initial_questions.crypto", {
      returnObjects: true,
    }) as string[]) || [],
  lover:
    (t("chat.initial_questions.lover", { returnObjects: true }) as string[]) ||
    [],
});

// Keep the original for backward compatibility
export const initialQuestions = {
  general: [
    "Unveil the horoscope for a soul born on [Birth Date] [Birth Time] [Birth Place]",
    "What Astro prophecy do the stars offer about [TOKEN] — will it rise or fall over the next 7 days?",
    "Destiny of my love life this year, for one born on [Birth Date] [Birth Time] [Birth Place]",
    "Reveal the cryptocurrency that aligns with my [Birth Date]",
  ],
  crypto: [
    "What does the cosmic energy reveal about [TOKEN] price movement in the next lunar cycle?",
    "Which cryptocurrency will shine brightest under the stars of [Month/Year]?",
    "Divine the fortune of my crypto portfolio born on [Date] with holdings in [TOKEN]",
    "What celestial signs indicate the best time to invest in [TOKEN]?",
  ],
  lover: [
    "Will love find me under the stars of [Birth Date] [Birth Time] [Birth Place]?",
    "Cosmic compatibility reading for souls born on [Your Birth Date] and [Partner Birth Date]",
    "What does Venus whisper about my romantic destiny this [Season/Month]?",
    "Divine guidance for healing a heart born on [Birth Date] after recent heartbreak",
  ],
};

export const getCategories = (t: TFunction) => [
  {
    id: "general",
    subtitle: t("common.chat"),
    title: t("common.general"),
    icon: Cat,
    color: "bg-red-700",
    hoverColor: "hover:from-red-700 hover:to-red-800",
  },
  {
    id: "crypto",
    subtitle: t("common.chat"),
    title: t("common.crypto"),
    icon: DollarSign,
    color: "bg-red-700",
    hoverColor: "hover:from-red-700 hover:to-red-800",
  },
  {
    id: "lover",
    subtitle: t("common.chat"),
    title: t("common.lover"),
    icon: Heart,
    color: "bg-red-700",
    hoverColor: "hover:from-red-700 hover:to-red-800",
  },
  {
    id: "tarot",
    subtitle: t("common.card"),
    title: "Tarot",
    icon: Diamond,
    color: "bg-white/50",
    hoverColor: "hover:from-red-700 hover:to-red-800",
    comingSoon: true,
  },
];

// Keep the original for backward compatibility
export const categories = [
  {
    id: "general",
    subtitle: "Chat",
    title: "General",
    icon: Cat,
    color: "bg-red-700",
    hoverColor: "hover:from-red-700 hover:to-red-800",
  },
  {
    id: "crypto",
    subtitle: "Chat",
    title: "Crypto",
    icon: DollarSign,
    color: "bg-red-700",
    hoverColor: "hover:from-red-700 hover:to-red-800",
  },
  {
    id: "lover",
    subtitle: "Chat",
    title: "Lover",
    icon: Heart,
    color: "bg-red-700",
    hoverColor: "hover:from-red-700 hover:to-red-800",
  },
  {
    id: "tarot",
    subtitle: "Card",
    title: "Tarot",
    icon: Diamond,
    color: "bg-white/50",
    hoverColor: "hover:from-red-700 hover:to-red-800",
    comingSoon: true,
  },
];

export type CategoryId = keyof typeof initialQuestions;
