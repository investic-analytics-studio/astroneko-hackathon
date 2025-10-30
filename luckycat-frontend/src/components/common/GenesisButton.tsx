import { useNavigate } from "@tanstack/react-router";
import { useLanguage } from "@/hooks";

export function GenesisButton() {
  const navigate = useNavigate();
  const { currentLanguage } = useLanguage();

  const handleClick = () => {
    navigate({ to: "/$lng/genesis", params: { lng: currentLanguage } });
  };

  return (
    <div className="relative">
      {/* Sparkle animation background */}
      <div className="absolute inset-0 rounded-full bg-gradient-to-r from-purple-400 via-pink-500 to-red-500 opacity-75 blur-sm animate-pulse"></div>

      <button
        onClick={handleClick}
        className="relative bg-gradient-to-r from-yellow-400 via-pink-500 to-purple-600 text-white text-[14px] md:text-[16px] rounded-full pr-4 md:pr-2 py-2 md:py-2 gap-2 md:gap-4 flex items-center justify-center
        hover:from-yellow-300 hover:via-pink-400 hover:to-purple-500 hover:scale-105 hover:shadow-lg hover:shadow-purple-500/50
        transition-all duration-300 cursor-pointer focus:none ring-0 outline-none border-0
        animate-shimmer bg-[length:200%_100%] font-semibold
        shadow-lg shadow-purple-500/25
        before:absolute before:inset-0 before:rounded-full before:bg-gradient-to-r before:from-transparent before:via-white/20 before:to-transparent before:animate-shine
        "
        style={{
          background:
            "linear-gradient(45deg, #fbbf24, #ec4899, #8b5cf6, #fbbf24)",
          backgroundSize: "200% 200%",
          animation:
            "gradient-shift 3s ease infinite, sparkle 2s ease-in-out infinite alternate",
        }}
      >
        <span className="relative z-10 flex items-center gap-2 md:gap-4">
          ✨ Astro Genesis{" "}
          <span className="hidden md:block bg-gradient-to-r from-red-500 to-pink-600 text-[12px] md:text-[16px] text-white rounded-full px-2 py-1 text-xs shadow-md">
            New
          </span>
        </span>
      </button>

      {/* Additional sparkle effects */}
      <div className="absolute -top-1 -right-1 w-3 h-3 bg-yellow-400 rounded-full animate-ping opacity-75"></div>
      <div
        className="absolute -bottom-1 -left-1 w-2 h-2 bg-pink-400 rounded-full animate-ping opacity-75"
        style={{ animationDelay: "0.5s" }}
      ></div>
    </div>
  );
}
