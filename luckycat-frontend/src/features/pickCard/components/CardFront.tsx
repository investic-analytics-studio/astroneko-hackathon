import { motion } from "framer-motion";
import { Card } from "@/types/pickCard";
import { OptimizedImage } from "@/components/ui/OptimizedImage";
import { useMobileDetection } from "@/hooks/common/useDeviceDetection";

interface CardFrontProps {
  card: Card;
}

export const CardFront = ({ card }: CardFrontProps) => {
  const isMobile = useMobileDetection();

  return (
    <motion.div
      className="absolute inset-0 w-full h-full backface-hidden will-change-transform card-front-hidden"
      style={{
        backfaceVisibility: "hidden",
        WebkitBackfaceVisibility: "hidden",
        transform: "rotateY(180deg)",
        transformStyle: "preserve-3d",
        // Ensure proper 3D rendering on mobile
        WebkitTransform: "rotateY(180deg)",
        WebkitTransformStyle: "preserve-3d",
        // Additional mobile-specific styles
        ...(isMobile && {
          transform: "rotateY(180deg) translateZ(0)",
          WebkitTransform: "rotateY(180deg) translateZ(0)",
        }),
      }}
    >
      <div className="relative w-full h-full rounded-xl border-[3px] border-amber-500/30 overflow-hidden shadow-[0_0_15px_rgba(245,158,11,0.3)]">
        {/* Decorative corners */}
        <div className="absolute top-0 left-0 w-6 h-6 border-t-2 border-l-2 border-amber-500/50 rounded-tl-lg z-10"></div>
        <div className="absolute top-0 right-0 w-6 h-6 border-t-2 border-r-2 border-amber-500/50 rounded-tr-lg z-10"></div>
        <div className="absolute bottom-0 left-0 w-6 h-6 border-b-2 border-l-2 border-amber-500/50 rounded-bl-lg z-10"></div>
        <div className="absolute bottom-0 right-0 w-6 h-6 border-b-2 border-r-2 border-amber-500/50 rounded-br-lg z-10"></div>
        {/* Inner glow */}
        <div className="absolute inset-0 rounded-xl shadow-[inset_0_0_15px_rgba(245,158,11,0.2)] z-10"></div>
        <OptimizedImage
          src={card.image}
          alt={card.name}
          className="w-full h-full object-cover"
          loading="lazy"
          priority={false}
        />
      </div>
    </motion.div>
  );
};
