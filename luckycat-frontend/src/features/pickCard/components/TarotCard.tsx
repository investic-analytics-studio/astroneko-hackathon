import { motion } from "framer-motion";

import { MOTION_EASE } from "@/constants/motion";
import { Card } from "@/types/pickCard";
import { CardFront } from "./CardFront";
import { CardBack } from "./CardBack";

interface TarotCardProps {
  card: Card;
  index: number;
  selectedIndex: number | null;
  onClick: () => void;
}

export const TarotCard = ({
  card,
  index,
  selectedIndex,
  onClick,
}: TarotCardProps) => (
  <motion.div
    key={card.name}
    className="cursor-pointer relative w-full aspect-[2/3] [perspective:1000px] group"
    style={{
      perspective: "1000px",
      WebkitPerspective: "1000px",
    }}
    initial={{ opacity: 0, y: 20, scale: 0.9 }}
    animate={{ opacity: 1, y: 0, scale: 1 }}
    transition={{
      duration: 0.6,
      delay: index * 0.15,
      ease: MOTION_EASE,
    }}
    whileHover={
      selectedIndex === null
        ? {
            scale: 1.05,
            y: -10,
            transition: {
              duration: 0.3,
              ease: MOTION_EASE,
            },
          }
        : {}
    }
    onClick={() => selectedIndex === null && onClick()}
  >
    <motion.div
      className="relative w-full h-full [transform-style:preserve-3d] transition-all duration-500"
      style={{
        transformStyle: "preserve-3d",
        WebkitTransformStyle: "preserve-3d",
      }}
      animate={
        selectedIndex === index
          ? {
              rotateY: 180,
              scale: 1.1,
              transition: {
                duration: 0.6,
                ease: MOTION_EASE,
              },
          }
        : selectedIndex !== null
        ? {
            scale: 0.95,
            opacity: 0.5,
            filter: "blur(2px)",
            transition: {
              duration: 0.4,
              ease: MOTION_EASE,
            },
          }
        : {}
      }
    >
      <CardBack />
      <CardFront card={card} />
    </motion.div>
  </motion.div>
);
