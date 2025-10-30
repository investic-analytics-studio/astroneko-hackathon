import { motion } from "framer-motion";
import { getImagePath } from "@/lib/utils";
import { useMobileDetection } from "@/hooks/common/useDeviceDetection";

export const CardBack = () => {
  const isMobile = useMobileDetection();

  return (
    <div
      className="relative w-full h-full rounded-[10px] overflow-hidden shadow-[0_0px_8px_1px_rgba(0,0,0,0.8)] card-back-visible"
      style={{
        backfaceVisibility: "hidden",
        WebkitBackfaceVisibility: "hidden",
        transform: "rotateY(0deg)",
        transformStyle: "preserve-3d",
        WebkitTransform: "rotateY(0deg)",
        WebkitTransformStyle: "preserve-3d",
        // Additional mobile-specific styles
        ...(isMobile && {
          transform: "rotateY(0deg) translateZ(0)",
          WebkitTransform: "rotateY(0deg) translateZ(0)",
        }),
      }}
    >
      {/* Card background with full image */}
      <img
        src={getImagePath("/cards/card-back-3.webp")}
        alt="Card back"
        className="absolute inset-0 w-full h-full object-fill"
      />

      {/* Animated stars pattern */}
      <motion.div
        className="absolute inset-0"
        animate={{
          opacity: [0.2, 0.4, 0.2],
        }}
        transition={{
          duration: 2,
          repeat: Infinity,
          repeatType: "reverse",
        }}
      />
    </div>
  );
};
