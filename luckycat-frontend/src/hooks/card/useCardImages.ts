import { useMemo } from "react";
import { allCards } from "@/constants/pickCard";
import { useImagePreloader } from "@/hooks";
import { getImagePath } from "@/lib/utils";

const getCardImagePaths = (): string[] =>
  allCards.map((card) => card.image.replace(getImagePath(""), ""));

export const useCardImages = () => {
  const images = useMemo(
    () => [
      "cards/card-back-3.webp",
      "cards/card-back.webp",
      ...getCardImagePaths(),
    ],
    []
  );

  useImagePreloader({ images, priority: false });
};
