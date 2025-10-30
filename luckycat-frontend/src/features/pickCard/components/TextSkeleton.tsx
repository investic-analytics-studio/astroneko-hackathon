import { motion } from "framer-motion";

export const TextSkeleton = () => (
  <div className="space-y-6 w-full mt-4 pl-8">
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-3/4"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-4/6"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.1,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-5/6"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.2,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-4/5"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.3,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-5/6"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.2,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-4/5"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.3,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-5/6"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.2,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-4/5"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.3,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="hidden 2xl:block h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-5/6"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.2,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
    <motion.div
      className="hidden 2xl:block h-4 2xl:h-10 bg-gradient-to-r from-[#FFFFFF]/20 via-[#FFFFFF]/30 to-[#FFFFFF]/20 rounded-full w-4/5"
      animate={{
        opacity: [0.5, 1, 0.5],
      }}
      transition={{
        duration: 2,
        delay: 0.3,
        repeat: Infinity,
        ease: "easeInOut",
      }}
    />
  </div>
);
