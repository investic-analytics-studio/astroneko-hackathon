import { useCallback } from "react";
import { useCardSelection } from "./useCardSelection";
import { useFortuneTelling } from "./useFortuneTelling";
import { useCardImages } from "./useCardImages";
import { sleep } from "@/lib/utils";
import type { Card } from "@/types/pickCard";
import { track } from "@/lib/amplitude";

const TIMING = {
  PRE_REVEAL_WAIT: 600,
  PRE_RESULT_WAIT: 500,
  RESET_TRANSITION: 400,
} as const;

export interface UsePickCardControllerReturn {
  selectedCards: Card[];
  finalCard: Card | null;
  isSelecting: boolean;
  selectedIndex: number | null;
  fortuneReply: string;
  isLoading: boolean;
  isInitialLoading: boolean;
  handleCardClick: (card: Card, index: number) => Promise<void>;
  handleReset: () => void;
}

export const usePickCardController = (): UsePickCardControllerReturn => {
  const {
    selectedCards,
    finalCard,
    selectedIndex,
    isSelecting,
    isInitialLoading,
    selectCard,
    reset: resetSelection,
  } = useCardSelection();

  const { fortuneReply, isLoading, getFortune, clearReply } =
    useFortuneTelling();

  useCardImages();

  const handleCardClick = useCallback(
    async (card: Card, index: number) => {
      if (isLoading) return;

      track('card selected', {
        card_name: card.name,
        card_index: index,
        card_prompt: card.prompt,
      });

      selectCard(card, index);
      await sleep(TIMING.PRE_REVEAL_WAIT);
      await getFortune(card.prompt);
      await sleep(TIMING.PRE_RESULT_WAIT);

      track('fortune revealed', {
        card_name: card.name,
        has_fortune: !!fortuneReply,
      });
    },
    [isLoading, selectCard, getFortune, fortuneReply]
  );

  const handleReset = useCallback(() => {
    if (isLoading) return;

    track('card selection reset', {
      had_selection: !!finalCard,
    });

    setTimeout(() => {
      resetSelection();
      clearReply();
    }, TIMING.RESET_TRANSITION);
  }, [isLoading, resetSelection, clearReply, finalCard]);

  return {
    selectedCards,
    finalCard,
    isSelecting,
    selectedIndex,
    fortuneReply,
    isLoading,
    isInitialLoading,
    handleCardClick,
    handleReset,
  };
};
