import { useState, useCallback, useEffect } from "react";
import { allCards } from "@/constants/pickCard";
import type { Card } from "@/types/pickCard";

const CARDS_TO_DISPLAY = 6;
const INITIAL_LOADING_DELAY = 100;

const getShuffledCards = (count: number): Card[] => {
  const shuffled = [...allCards].sort(() => Math.random() - 0.5);
  return shuffled.slice(0, count);
};

export const useCardSelection = () => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);
  const [finalCard, setFinalCard] = useState<Card | null>(null);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const [isSelecting, setIsSelecting] = useState(true);
  const [isInitialLoading, setIsInitialLoading] = useState(true);

  const initializeCards = useCallback(() => {
    setIsInitialLoading(true);
    setSelectedCards(getShuffledCards(CARDS_TO_DISPLAY));
    setTimeout(() => setIsInitialLoading(false), INITIAL_LOADING_DELAY);
  }, []);

  const selectCard = useCallback((card: Card, index: number) => {
    setSelectedIndex(index);
    setFinalCard(card);
    setIsSelecting(false);
  }, []);

  const reset = useCallback(() => {
    setIsSelecting(true);
    setFinalCard(null);
    setSelectedIndex(null);
    initializeCards();
  }, [initializeCards]);

  useEffect(() => {
    initializeCards();
  }, [initializeCards]);

  return {
    selectedCards,
    finalCard,
    selectedIndex,
    isSelecting,
    isInitialLoading,
    selectCard,
    reset,
  };
};
