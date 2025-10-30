import { Clock } from "lucide-react";

export const ComingSoonBadge = () => (
  <div
    className="absolute top-4 right-4 z-20 opacity-0 translate-y-2 animate-fade-in"
    style={{ animationDelay: "0.5s" }}
  >
    <div className="bg-yellow-500 rounded-2xl px-3 py-2 shadow-lg">
      <div className="flex items-center justify-center space-x-2">
        <div className="animate-bounce-slow">
          <Clock className="w-3 h-3 sm:w-4 sm:h-4 text-white" />
        </div>
        <span className="text-white font-bold text-xs sm:text-sm font-press-start">
          Coming Soon
        </span>
      </div>
    </div>
  </div>
);
