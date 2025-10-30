import { memo } from "react";

import { motion } from "framer-motion";

import { MOTION_EASE, MOTION_EASE_OUT } from "@/constants/motion";
import { FortuneDisplay } from "@/features/pickCard";
import { getImagePath } from "@/lib/utils";
import type { Card } from "@/types/pickCard";

interface RevealViewProps {
  card: Card | null;
  fortuneReply: string;
  isLoading: boolean;
  onReset: () => void;
}

const RevealViewComponent = ({
  card,
  fortuneReply,
  isLoading,
  onReset,
}: RevealViewProps) => {
  if (!card) {
    return null;
  }

  return (
    <motion.div
      key="pickcard-reveal"
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      exit={{ opacity: 0, scale: 0.9 }}
      transition={{ duration: 0.5, ease: MOTION_EASE }}
    >
      <div className="flex min-h-screen flex-col items-center justify-center overflow-y-hidden pt-10 md:min-h-0 md:items-start md:justify-start md:pt-60 xl:grid xl:grid-cols-3 xl:pt-10">
        <div className="w-full max-w-sm xl:col-span-1 xl:h-[700px] xl:p-10">
          <motion.div
            className="group flex h-full w-full items-center justify-center rounded-xl"
            initial={{ opacity: 0, rotateY: 180 }}
            animate={{ opacity: 1, rotateY: 0 }}
            transition={{ duration: 0.6, ease: MOTION_EASE_OUT }}
            style={{ willChange: "transform", transformStyle: "preserve-3d" }}
          >
            <img
              src={getImagePath(card.image)}
              alt={card.name}
              className="block h-[350px] w-auto rounded-xl shadow-[0_0px_8px_1px_rgba(0,0,0,0.8)] md:h-[600px] xl:h-full xl:w-full"
            />
          </motion.div>
        </div>

        <div className="xl:col-span-2 xl:h-[700px] xl:p-10">
          <motion.div
            className="min-w-0 h-full lg:w-full 2xl:w-full"
            initial={{ opacity: 0, x: 30 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: 0.4, duration: 0.6 }}
          >
            <motion.div
              className="relative h-full w-full rounded-xl border border-white/20 bg-black/70 backdrop-blur-md xl:min-h-[550px] xl:max-h-[550px] xl:overflow-y-auto"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.5, delay: 0.6 }}
            >
              <FortuneDisplay
                isLoading={isLoading}
                fortuneReply={fortuneReply}
              />
            </motion.div>

            <motion.button
              onClick={onReset}
              disabled={isLoading}
              className={`font-press-start relative z-20 mt-6 flex h-[50px] w-auto items-center justify-center gap-2 rounded-full px-8 py-3 text-[12px] font-semibold transition-all duration-300 shadow-lg backdrop-blur-sm focus:outline-none md:text-[16px] 2xl:mt-10 2xl:h-[120px] 2xl:text-[30px] ${
                isLoading
                  ? "cursor-not-allowed bg-[#F7C36D]/80 text-black/80"
                  : "bg-[#F7C36D] text-black hover:scale-105 hover:shadow-[0_0_20px_var(--bg-gradient-1)]"
              }`}
              whileHover={!isLoading ? { scale: 1.02 } : {}}
              whileTap={!isLoading ? { scale: 0.95 } : {}}
            >
              {isLoading ? "Loading..." : "Draw A New Astro Card"}
            </motion.button>
          </motion.div>
        </div>
      </div>
    </motion.div>
  );
};

export const RevealView = memo(RevealViewComponent);

RevealView.displayName = "RevealView";
