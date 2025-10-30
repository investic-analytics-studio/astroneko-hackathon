import { type ReactNode } from "react";

import { AnimatePresence, motion } from "framer-motion";

import { MOTION_EASE } from "@/constants/motion";
import { CopyButton } from "@/components/ui/copy-button";
import { useFormattedFortune } from "./fortuneText";
import { TextSkeleton } from "./TextSkeleton";

interface FortuneDisplayProps {
  isLoading: boolean;
  fortuneReply: string;
}

interface FortuneDisplayHeaderProps {
  showCopyButton: boolean;
  fortuneReply: string;
}

export const FortuneDisplayHeader = ({
  showCopyButton,
  fortuneReply,
}: FortuneDisplayHeaderProps) => (
  <motion.div
    className="flex items-center justify-between w-full gap-2 mb-4 px-10 py-6 border-b border-white/15 bg-white/10"
    initial={{ opacity: 0, y: -10 }}
    animate={{ opacity: 1, y: 0 }}
    transition={{ duration: 0.4, delay: 0.2, ease: MOTION_EASE }}
  >
    <motion.h2
      className="text-[14px] xl:text-[20px] 2xl:text-[32px] font-semibold text-[#F7C36D] font-press-start text-center"
      initial={{ opacity: 0, y: -5 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: 0.3, ease: MOTION_EASE }}
    >
      Astro Feline Whispers
    </motion.h2>
    {showCopyButton && (
      <CopyButton
        textToCopy={fortuneReply}
        successMessage="Fortune copied!"
        variant="ghost"
        size="sm"
        className="text-[12px] 2xl:text-sm p-2 rounded-md border border-white/20 bg-white/5 hover:bg-white/10 hover:border-white/50 text-white"
        showLabel={false}
      />
    )}
  </motion.div>
);

interface FortuneContentProps {
  isLoading: boolean;
  formattedFortune: ReactNode[];
}

export const FortuneEmptyState = () => (
  <motion.div
    key="fortune-empty"
    initial={{ opacity: 0, y: 10 }}
    animate={{ opacity: 1, y: 0 }}
    exit={{ opacity: 0, y: 10 }}
    transition={{ duration: 0.3, ease: MOTION_EASE }}
    className="flex h-full items-center justify-center"
  >
    <p className="text-center text-sm text-white/70 sm:text-base">
      Ask Astro Feline for guidance to reveal your fortune.
    </p>
  </motion.div>
);

export const FortuneContent = ({
  isLoading,
  formattedFortune,
}: FortuneContentProps) => {
  let content: ReactNode = null;

  if (isLoading) {
    content = <TextSkeleton key="fortune-loading" />;
  } else if (formattedFortune.length === 0) {
    content = <FortuneEmptyState />;
  } else {
    content = (
      <motion.div
        key="fortune-content"
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -15 }}
        transition={{ duration: 0.5, delay: 0.2, ease: MOTION_EASE }}
        className="relative h-full flex flex-col"
      >
        <motion.div
          className="h-full px-4 sm:px-6 2xl:px-0 py-3 sm:py-4 space-y-4 w-full text-lg sm:text-base xl:text-lg 2xl:text-[40px]"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.4, delay: 0.3 }}
        >
          <div className="mx-auto text-left pb-4">{formattedFortune}</div>
        </motion.div>
      </motion.div>
    );
  }

  return <AnimatePresence mode="wait">{content}</AnimatePresence>;
};

/**
 * Component to display generated fortune text with rich formatting and animations.
 */
export const FortuneDisplay = ({
  isLoading,
  fortuneReply,
}: FortuneDisplayProps) => {
  const formattedFortune = useFormattedFortune(fortuneReply);
  const hasFortuneContent = formattedFortune.length > 0;
  const showCopyButton = hasFortuneContent && !isLoading;

  return (
    <div className="relative w-full h-full p-0">
      <motion.div
        className="relative h-full flex flex-col items-center space-y-2 sm:space-y-4 w-full max-w-[95vw] sm:max-w-4xl 2xl:max-w-full mx-auto"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: 20 }}
        transition={{ duration: 0.4, ease: MOTION_EASE }}
      >
        <FortuneDisplayHeader
          showCopyButton={showCopyButton}
          fortuneReply={fortuneReply}
        />

        <div className="w-full max-w-full sm:max-w-5xl 2xl:max-w-full mx-auto px-2 sm:px-4 2xl:px-0 flex-1 min-h-0">
          <FortuneContent
            isLoading={isLoading}
            formattedFortune={formattedFortune}
          />
        </div>

        <motion.div
          className="w-32 h-px mt-6"
          initial={{ opacity: 0, scaleX: 0.8 }}
          animate={{ opacity: 1, scaleX: 1 }}
          exit={{ opacity: 0, scaleX: 0.8 }}
          transition={{ duration: 0.4, delay: 0.5 }}
        >
          <div className="w-full h-full bg-gradient-to-r from-[var(--bg-gradient-1)] to-[var(--bg-gradient-2)]" />
        </motion.div>
      </motion.div>
    </div>
  );
};

FortuneDisplay.displayName = "FortuneDisplay";
