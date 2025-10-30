import { useMemo, type ReactNode } from "react";

import { motion } from "framer-motion";

import { MOTION_EASE, MOTION_SPRING_EASE } from "@/constants/motion";

const TOKEN_SPLIT_REGEX = /(\*\*.*?\*\*|\*[^*]+\*|\n+|🔮|✨|---|-|###)/;
const BASE_DELAY = 0.15;
const STAGGER_DELAY = 0.02;

const getDelay = (index: number) => BASE_DELAY + index * STAGGER_DELAY;

export const normalizeFortuneText = (text: string) => text.replace(/\\n/g, "\n");

const renderFormattedToken = (token: string, index: number): ReactNode => {
  if (!token) {
    return null;
  }

  if (token === "###") {
    return null;
  }

  if (/\n+/.test(token)) {
    return (
      <motion.div
        key={`break-${index}`}
        initial={{ opacity: 0, height: 0 }}
        animate={{ opacity: 1, height: 8 }}
        transition={{ delay: getDelay(index), duration: 0.4, ease: MOTION_EASE }}
      />
    );
  }

  const trimmedToken = token.trim();

  if (!trimmedToken) {
    return null;
  }

  if (token.startsWith("**") && token.endsWith("**")) {
    const emphasizedText = token.slice(2, -2).replace(/_/g, " ");

    return (
      <motion.span
        key={`bold-${index}`}
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: getDelay(index), duration: 0.6, ease: MOTION_EASE }}
        className="inline-block font-medium bg-gradient-to-br from-[var(--bg-gradient-1)] via-[var(--bg-gradient-2)] to-[var(--bg-gradient-2)] text-black px-3 py-1 rounded-lg mx-1 border border-[var(--bg-gradient-1)] shadow-[0_0_10px_var(--bg-gradient-1)] backdrop-blur-sm hover:scale-105 hover:[var(--bg-gradient-1)] transition-all duration-300"
      >
        {emphasizedText}
      </motion.span>
    );
  }

  if (token.startsWith("*") && token.endsWith("*") && !token.startsWith("**")) {
    const emphasizedText = token.slice(1, -1);

    return (
      <motion.span
        key={`italic-${index}`}
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: getDelay(index), duration: 0.5, ease: MOTION_EASE }}
        className="inline-block italic font-medium text-cyan-300/90 px-1 shadow-[0_0_8px_rgba(103,232,249,0.3)] hover:text-cyan-200 transition-colors duration-200"
      >
        {emphasizedText}
      </motion.span>
    );
  }

  if (token === "🔮" || token === "✨") {
    return (
      <motion.span
        key={`emoji-${index}`}
        initial={{ scale: 0.5, opacity: 0, rotate: -15 }}
        animate={{ scale: 1, opacity: 1, rotate: 0 }}
        transition={{ delay: getDelay(index), duration: 0.6, ease: MOTION_SPRING_EASE }}
        whileHover={{ scale: 1.2, rotate: 15, transition: { duration: 0.3, ease: "easeOut" } }}
        className="inline-block mx-1 cursor-default"
      >
        {token}
      </motion.span>
    );
  }

  if (token === "---") {
    return null;
  }

  if (token === "-") {
    return (
      <motion.span
        key={`dash-${index}`}
        initial={{ scaleX: 0, opacity: 0 }}
        animate={{ scaleX: 1, opacity: 1 }}
        transition={{ delay: getDelay(index), duration: 0.4, ease: MOTION_EASE }}
        className="inline-block mx-1 w-3 h-[2px] bg-gradient-to-r from-[var(--bg-gradient-1)] to-[var(--bg-gradient-2)] rounded-full align-middle translate-y-[-2px]"
      />
    );
  }

  return (
    <motion.span
      key={`text-${index}`}
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: getDelay(index), duration: 0.5, ease: MOTION_EASE }}
      className="text-gray-100/90 leading-relaxed tracking-wide whitespace-pre-wrap"
    >
      {token}
    </motion.span>
  );
};

export const formatFortuneText = (text: string): ReactNode[] => {
  if (!text) {
    return [];
  }

  const normalizedText = normalizeFortuneText(text);
  const tokens = normalizedText.split(TOKEN_SPLIT_REGEX);

  return tokens.reduce<ReactNode[]>((acc, token, index) => {
    const renderedToken = renderFormattedToken(token, index);

    if (renderedToken) {
      acc.push(renderedToken);
    }

    return acc;
  }, []);
};

export const useFormattedFortune = (fortuneReply: string) =>
  useMemo(() => formatFortuneText(fortuneReply), [fortuneReply]);
