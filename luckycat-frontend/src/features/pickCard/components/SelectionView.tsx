import { motion } from "framer-motion";
import { useTranslation } from "react-i18next";

import {
  MOTION_EASE_IN_OUT,
  MOTION_EASE_OUT,
  MOTION_STAGGER_SM,
} from "@/constants/motion";
import type { Card } from "@/types/pickCard";

import { TarotCard } from "./TarotCard";

interface SelectionViewProps {
  isInitialLoading: boolean;
  isActive: boolean;
  selectedCards: Card[];
  selectedIndex: number | null;
  onCardSelect: (card: Card, index: number) => void;
}

export const SkeletonCard = ({ index }: { index: number }) => (
  <motion.div
    key={`loading-${index}`}
    className="relative w-full aspect-[2/3]"
    initial={{ opacity: 0, scale: 0.9 }}
    animate={{
      opacity: [0.4, 0.7, 0.4],
      scale: [0.98, 1, 0.98],
    }}
    transition={{
      duration: 2,
      repeat: Infinity,
      repeatType: "reverse",
      delay: index * 0.1,
      ease: MOTION_EASE_IN_OUT,
    }}
  >
    <div className="h-full w-full rounded-xl border-2 border-purple-500/20 bg-gradient-to-br from-purple-500/10 to-pink-500/5" />
  </motion.div>
);

export const SelectionView = ({
  isInitialLoading,
  isActive,
  selectedCards,
  selectedIndex,
  onCardSelect,
}: SelectionViewProps) => {
  const { t } = useTranslation();

  return (
    <>
      <motion.div
        key="pickcard-header"
        className="relative z-10 items-center gap-0 px-4 pt-40 md:px-10 md:pt-20 lg:pt-40 xl:pt-24"
        initial={{ opacity: 0, y: -20 }}
        animate={{
          opacity: 1,
          y: 0,
          scale: isActive ? 1 : 0.9,
        }}
        transition={{ duration: 0.5, ease: MOTION_EASE_OUT }}
      >
        <motion.h1
          className="font-press-start text-center text-[14px] font-semibold leading-[20px] text-white md:text-[24px] md:leading-[40px] lg:text-[30px] lg:leading-[40px] xl:px-[100px] 2xl:text-[40px] 2xl:leading-[60px]"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.2 }}
        >
          {t("pickcard.universe_question")}
        </motion.h1>
        <motion.h3
          className="font-press-start text-center text-[18px] font-semibold text-[#F7C36D] md:text-[24px] md:leading-tight lg:text-[30px] xl:mt-10 xl:text-[30px] 2xl:mt-20 2xl:text-[50px]"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.8, delay: 0.4 }}
        >
          {t("pickcard.focus_choose")}
        </motion.h3>
      </motion.div>

      <motion.div
        key="pickcard-grid"
        className="relative z-10 mx-auto grid w-full max-w-[1200px] grid-cols-3 justify-center gap-4 px-4 md:max-w-[600px] md:gap-4 lg:grid-cols-3 xl:grid-cols-6 xl:max-w-[1200px] 2xl:grid-cols-3"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, scale: 0.95, transition: { duration: 0.3 } }}
        transition={{
          duration: 0.5,
          ease: MOTION_EASE_OUT,
          staggerChildren: MOTION_STAGGER_SM,
        }}
      >
        {isInitialLoading
          ? Array.from({ length: 6 }, (_, index) => (
              <SkeletonCard key={index} index={index} />
            ))
          : selectedCards.map((card, index) => (
              <TarotCard
                key={card.name}
                card={card}
                index={index}
                selectedIndex={selectedIndex}
                onClick={() => onCardSelect(card, index)}
              />
            ))}
      </motion.div>
    </>
  );
};
