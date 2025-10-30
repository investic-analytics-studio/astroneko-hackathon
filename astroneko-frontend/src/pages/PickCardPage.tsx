import { AnimatePresence } from "framer-motion";

import { RevealView } from "@/features/pickCard";
import { SelectionView } from "@/features/pickCard";
import { usePickCardController } from "@/hooks";

export default function PickCardTarot() {
  const {
    selectedCards,
    finalCard,
    isSelecting,
    selectedIndex,
    fortuneReply,
    isLoading,
    isInitialLoading,
    handleCardClick,
    handleReset,
  } = usePickCardController();

  return (
    <div className="flex flex-col items-center bg-[image:var(--bg-display-2)] bg-cover bg-center justify-start min-h-screen overflow-y-auto pt-[64px]">
      <AnimatePresence mode="wait">
        {isSelecting ? (
          <SelectionView
            isInitialLoading={isInitialLoading}
            isActive={isSelecting}
            selectedCards={selectedCards}
            selectedIndex={selectedIndex}
            onCardSelect={handleCardClick}
          />
        ) : (
          <RevealView
            card={finalCard}
            fortuneReply={fortuneReply}
            isLoading={isLoading}
            onReset={handleReset}
          />
        )}
      </AnimatePresence>
    </div>
  );
}
