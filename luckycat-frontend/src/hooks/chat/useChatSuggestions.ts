import { useCallback, useMemo } from "react";
import { useTranslation } from "react-i18next";
import { getInitialQuestions, type CategoryId } from "@/constants/categories";
import { useChatStore } from "@/store/chatStore";

const DEFAULT_CATEGORY = "general";

export const useChatSuggestions = (category: string) => {
  const { t } = useTranslation();
  const { getCurrentCategoryState, setShowQuestions } = useChatStore();
  const { showQuestions } = getCurrentCategoryState(category);

  const questions = useMemo(() => {
    const translatedQuestions = getInitialQuestions(t);
    const categoryQuestions =
      translatedQuestions[category as CategoryId] ||
      translatedQuestions[DEFAULT_CATEGORY];
    return Array.isArray(categoryQuestions) ? categoryQuestions : [];
  }, [category, t]);

  const toggle = useCallback(() => {
    setShowQuestions(category, !showQuestions);
  }, [category, setShowQuestions, showQuestions]);

  const hide = useCallback(() => {
    setShowQuestions(category, false);
  }, [category, setShowQuestions]);

  return {
    questions,
    showQuestions,
    toggle,
    hide,
  };
};
