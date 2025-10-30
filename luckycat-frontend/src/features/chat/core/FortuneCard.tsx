import { motion } from "framer-motion";
import { memo } from "react";
import ReactMarkdown from "react-markdown";
import rehypeRaw from "rehype-raw";
import remarkGfm from "remark-gfm";
import type { Components } from "react-markdown";

import { MOTION_EASE } from "@/constants/motion";
import { OptimizedImage } from "@/components/ui/OptimizedImage";
import { normalizeFortuneText } from "@/features/pickCard/components/fortuneText";

interface FortuneCardProps {
  card: string;
  meaning: string;
  text: string;
}

export const FormatText = (text: string) => {
  const normalizedText = normalizeFortuneText(text);

  // Custom components for react-markdown to maintain animations and styling
  const components: Components = {
    // Bold text with gradient background
    strong: ({ children }) => (
      <motion.span
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.6,
          ease: MOTION_EASE,
        }}
        className="inline-block font-normal bg-gradient-to-br from-[var(--bg-gradient-1)] via-[var(--bg-gradient-2)] to-[var(--bg-gradient-2)] 
                   text-black text-[13px] sm:text-[14px] md:text-[15px] lg:text-[16px] xl:text-[15px] 2xl:text-[17px] px-2 rounded-sm mx-1 border border-[var(--bg-gradient-1)] 
                   hover:scale-105 hover:[var(--bg-gradient-1)] transition-all duration-300"
      >
        {children}
      </motion.span>
    ),

    // Italic text
    em: ({ children }) => (
      <motion.span
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.5,
          ease: MOTION_EASE,
        }}
        className="inline-block italic font-medium text-cyan-300/90 px-1
                   shadow-[0_0_8px_rgba(103,232,249,0.3)]
                   hover:text-cyan-200 transition-colors duration-200"
      >
        {children}
      </motion.span>
    ),

    // Paragraphs
    p: ({ children }) => (
      <motion.p
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.5,
          ease: MOTION_EASE,
        }}
        className="bg-transparent border-none text-white text-[15px] sm:text-[16px] md:text-[17px] lg:text-[18px] xl:text-[17px] 2xl:text-[19px] font-light xl:font-medium rounded-[14px] 2xl:rounded-[30px] mb-4"
      >
        {children}
      </motion.p>
    ),

    // Headers
    h1: ({ children }) => (
      <motion.h1
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.6,
          ease: MOTION_EASE,
        }}
        className="text-lg sm:text-xl md:text-2xl lg:text-3xl xl:text-2xl 2xl:text-4xl font-bold text-white mb-4"
      >
        {children}
      </motion.h1>
    ),

    h2: ({ children }) => (
      <motion.h2
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.6,
          ease: MOTION_EASE,
        }}
        className="text-base sm:text-lg md:text-xl lg:text-2xl xl:text-xl 2xl:text-3xl font-bold text-white mb-3"
      >
        {children}
      </motion.h2>
    ),

    h3: ({ children }) => (
      <motion.h3
        initial={{ opacity: 0, y: 15 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.6,
          ease: MOTION_EASE,
        }}
        className="text-sm sm:text-base md:text-lg lg:text-xl xl:text-lg 2xl:text-2xl font-semibold text-white mb-2"
      >
        {children}
      </motion.h3>
    ),

    // Lists
    ul: ({ children }) => (
      <motion.ul
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.5,
          ease: MOTION_EASE,
        }}
        className="list-none text-white space-y-3 mb-4 pl-0"
      >
        {children}
      </motion.ul>
    ),

    ol: ({ children }) => (
      <motion.ol
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 0.5,
          ease: MOTION_EASE,
        }}
        className="list-none text-white space-y-3 mb-4 pl-0"
      >
        {children}
      </motion.ol>
    ),

    li: ({ children }) => (
      <motion.li
        initial={{ opacity: 0, x: -10 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{
          duration: 0.4,
          ease: MOTION_EASE,
        }}
        className="flex items-start gap-3 text-white relative"
      >
        {/* Custom bullet with mystical glow */}
        <div
          className="flex-shrink-0 w-2 h-2 mt-2.5 rounded-full bg-gradient-to-br from-[var(--bg-gradient-1)] to-[var(--bg-gradient-2)] 
                       shadow-[0_0_8px_rgba(var(--bg-gradient-1-rgb),0.4)] 
                       before:content-[''] before:absolute before:w-3 before:h-3 before:rounded-full 
                       before:bg-gradient-to-br before:from-[rgba(var(--bg-gradient-1-rgb),0.3)] before:to-[rgba(var(--bg-gradient-1-rgb),0.2)] 
                       before:blur-sm before:-z-10 before:top-0.5 before:left-0.5"
        />

        {/* Content */}
        <div className="flex-1 min-w-0">{children}</div>
      </motion.li>
    ),

    // Blockquotes
    blockquote: ({ children }) => (
      <motion.blockquote
        initial={{ opacity: 0, x: -20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{
          duration: 0.6,
          ease: MOTION_EASE,
        }}
        className="border-l-4 border-[var(--bg-gradient-1)] pl-4 italic text-gray-300 mb-4"
      >
        {children}
      </motion.blockquote>
    ),

    // Code blocks
    code: ({ children, className }) => {
      const isInline = !className;
      if (isInline) {
        return (
          <motion.code
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{
              duration: 0.4,
              ease: MOTION_EASE,
            }}
            className="bg-gray-800 text-cyan-300 px-1 py-0.5 rounded text-sm font-mono"
          >
            {children}
          </motion.code>
        );
      }
      return (
        <motion.pre
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{
            duration: 0.5,
            ease: [0.4, 0.0, 0.2, 1],
          }}
          className="bg-gray-800 p-4 rounded-lg overflow-x-auto mb-4"
        >
          <code className={className}>{children}</code>
        </motion.pre>
      );
    },

    // Links
    a: ({ children, href }) => (
      <motion.a
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{
          duration: 0.4,
          ease: [0.4, 0.0, 0.2, 1],
        }}
        href={href}
        className="text-cyan-400 hover:text-cyan-300 underline transition-colors duration-200"
        target="_blank"
        rel="noopener noreferrer"
      >
        {children}
      </motion.a>
    ),

    // Horizontal rules
    hr: () => (
      <motion.hr
        initial={{ scaleX: 0, opacity: 0 }}
        animate={{ scaleX: 1, opacity: 1 }}
        transition={{
          duration: 0.5,
          ease: [0.4, 0.0, 0.2, 1],
        }}
        className="w-full h-[1.5px] bg-gradient-to-r from-transparent via-[var(--bg-gradient-1)] to-transparent rounded-full my-6"
      />
    ),
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.3 }}
      className="prose prose-invert max-w-none"
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeRaw]}
        components={components}
      >
        {normalizedText}
      </ReactMarkdown>
    </motion.div>
  );
};

export const FormatMeaning = (meaningText: string) => {
  // Split by commas and clean up each part
  const parts = meaningText
    .split(",")
    .map((part) => part.trim())
    .map((part) => part.replace(/^and /, "").trim())
    .map((part) => part.replace(/\.$/, "").trim())
    .map((part) => part.toUpperCase())
    .filter((part) => part.length > 0);

  return (
    <div className="flex flex-wrap items-center justify-center gap-3">
      {parts.map((part, index) => (
        <motion.div
          key={index}
          initial={{ opacity: 0, scale: 0.9, y: 15 }}
          animate={{ opacity: 1, scale: 1, y: 0 }}
          transition={{
            delay: 0.15 * index,
            duration: 0.5,
            ease: "easeOut",
          }}
          whileHover={{
            scale: 1.05,
            y: -5,
            transition: { duration: 0.3 },
          }}
          className="relative group"
        >
          {/* Background card with mystical border */}
          <div
            className="absolute inset-0 bg-gradient-to-br from-[rgba(var(--bg-gradient-1-rgb),0.4)] via-[rgba(var(--bg-gradient-1-rgb),0.3)] to-[rgba(var(--bg-gradient-1-rgb),0.4)] 
                         rounded-xl blur-md transform group-hover:blur-lg transition-all duration-500"
          />

          <div
            className="relative px-2 rounded-sm
                         bg-gradient-to-br from-[var(--bg-gradient-1)] to-[var(--bg-gradient-2)]
                         border border-[var(--bg-gradient-1)]
                         group-hover:[var(--bg-gradient-1)] transition-all duration-500
                         overflow-hidden"
          >
            {/* Decorative elements */}
            <div
              className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[1px] 
                           bg-gradient-to-r from-transparent via-orange-400/40 to-transparent"
            />
            <div
              className="absolute bottom-0 left-1/2 -translate-x-1/2 w-full h-[1px]
                           bg-gradient-to-r from-transparent via-amber-400/40 to-transparent"
            />

            {/* Star decorations */}
            <span className="absolute top-1 right-2 text-[10px] text-orange-300/70">
              ✦
            </span>
            <span className="absolute top-2 left-2 text-[8px] text-amber-300/60">
              ⋆
            </span>
            <span className="absolute bottom-2 right-2 text-[8px] text-yellow-300/60">
              ⋆
            </span>
            <span className="absolute bottom-1 left-2 text-[10px] text-orange-300/70">
              ✦
            </span>

            {/* Main text */}
            <div className="relative z-10 px-1">
              <span className="text-black font-normal tracking-[0.08em] text-[12px] sm:text-[13px] md:text-[14px] lg:text-[15px] xl:text-[14px] 2xl:text-[16px]">
                {part}
              </span>
            </div>

            {/* Mystical sparkle overlay */}
            <motion.div
              className="absolute inset-0 opacity-0 group-hover:opacity-100
                         bg-[radial-gradient(circle_at_50%_50%,rgba(249,115,22,0.2)_0%,transparent_60%)]
                         transition-opacity duration-700"
              animate={{
                scale: [1, 1.2, 1],
                opacity: [0, 0.4, 0],
              }}
              transition={{
                duration: 2,
                repeat: Infinity,
                repeatType: "reverse",
              }}
            />
          </div>
        </motion.div>
      ))}
    </div>
  );
};

export const FortuneCard = memo(({ card, meaning, text }: FortuneCardProps) => {
  // Validate card prop to prevent empty src
  const cardImageSrc =
    card && card.trim()
      ? `/cards/${card.toUpperCase().replace(/ /g, "_")}.webp`
      : "/cards/THE_FOOL.webp"; // Fallback to a default card

  return (
    <div className="grid grid-cols-1 md:grid-cols-[1fr_auto] gap-6 max-w-full mx-auto w-full">
      <div className="space-y-6 min-h-0">
        {/* Fortune Text */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="bg-transparent rounded-xl border-none"
        >
          {FormatText(text)}
        </motion.div>

        {/* Card Meaning */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3, duration: 0.5 }}
        >
          <div className="font-medium text-white flex items-center gap-2 mb-2 px-1">
            <span className="text-xl 2xl:text-[28px]">
              <OptimizedImage
                src="/icons/cat-foot.webp"
                alt="Crystal Ball"
                className="w-4 h-4"
                loading="eager"
                priority={true}
              />
            </span>
            <span className="text-sm sm:text-base md:text-lg lg:text-xl xl:text-lg 2xl:text-2xl tracking-wide">
              Astro Insight
            </span>
          </div>
          <div className="py-2">{FormatMeaning(meaning)}</div>
        </motion.div>
      </div>

      {/* Card Image */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ delay: 0.2, duration: 0.5 }}
        className="w-full md:w-[320px] 2xl:w-[500px] aspect-[2/3] 2xl:aspect-[1/2] rounded-xl overflow-hidden self-start sticky top-6"
      >
        <OptimizedImage
          src={cardImageSrc}
          alt={card || "Tarot Card"}
          className="w-full h-full object-cover rounded-xl border-none shadow-2xl transform transition-transform hover:scale-[1.02] duration-300"
          loading="eager"
          priority={true}
        />
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="absolute bottom-0 left-0 right-0 p-6 text-center font-medium z-20 bg-gradient-to-t from-indigo-950/90 via-indigo-950/50 to-transparent"
        >
          <motion.div
            className="relative inline-block"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.6 }}
          >
            <span className="absolute inset-0 bg-gradient-to-r from-yellow-200/0 via-yellow-200/10 to-yellow-200/0 blur-sm" />
            <span className="relative text-lg tracking-wider font-medium bg-gradient-to-br from-yellow-100 via-indigo-100 to-yellow-100 text-transparent bg-clip-text drop-shadow-[0_0_8px_rgba(234,179,8,0.3)]">
              {card}
            </span>
            <motion.div
              className="absolute -inset-1 bg-gradient-to-r from-yellow-200/0 via-yellow-200/10 to-yellow-200/0 rounded-lg z-[-1]"
              animate={{
                opacity: [0.3, 0.6, 0.3],
                scale: [1, 1.02, 1],
              }}
              transition={{
                duration: 2,
                repeat: Infinity,
                repeatType: "reverse",
              }}
            />
          </motion.div>
        </motion.div>
      </motion.div>
    </div>
  );
});
