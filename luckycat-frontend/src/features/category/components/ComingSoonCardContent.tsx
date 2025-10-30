import { LucideIcon, Clock } from "lucide-react";

interface ComingSoonCardContentProps {
  title: string;
  subtitle: string;
  icon: LucideIcon;
}

export const ComingSoonCardContent = ({
  title,
  subtitle,
  icon: IconComponent,
}: ComingSoonCardContentProps) => (
  <div className="relative z-10 flex flex-row md:flex-col justify-between h-full gap-4 md:gap-6 xl:gap-16">
    {/* Coming Soon Badge - Absolutely positioned in center */}
    <div
      className="absolute inset-0 flex items-center justify-center z-20 opacity-0 translate-y-2 animate-fade-in"
      style={{ animationDelay: "0.5s" }}
    >
      <div className="bg-yellow-500 rounded-2xl px-3 py-2 sm:px-4 sm:py-3 shadow-lg">
        <div className="flex items-center justify-center space-x-2 sm:space-x-3">
          <div className="animate-bounce-slow">
            <Clock className="w-3 h-3 sm:w-4 sm:h-4 md:w-5 md:h-5 text-white" />
          </div>
          <span className="text-white font-bold text-xs sm:text-sm md:text-base font-press-start">
            Coming Soon
          </span>
        </div>
      </div>
    </div>

    <div className="flex flex-col md:flex-row gap-2 sm:gap-4 md:gap-6 lg:gap-8 xl:gap-10 items-start justify-between">
      <div className="flex flex-col items-start justify-start">
        <p className="text-xs sm:text-sm md:text-base lg:text-lg font-medium opacity-90 leading-relaxed">
          {subtitle}
        </p>
        <h3 className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold font-press-start">
          {title}
        </h3>
      </div>
      <div className="bg-white/20 rounded-full p-2 transition-colors duration-300 animate-spin-slow">
        <Clock className="w-5 h-5 md:w-6 md:h-6" />
      </div>
    </div>

    <div className="flex items-center justify-center">
      <div className="bg-white/20 rounded-full p-4 sm:p-4 md:p-6 lg:p-8 xl:p-10 transition-all ease-in-out duration-300 group-hover:bg-white/30 group-hover:scale-105 group-hover:rotate-5">
        <IconComponent className="w-6 h-6 md:w-16 md:h-16 lg:w-20 lg:h-20 xl:w-20 xl:h-20" />
      </div>
    </div>
  </div>
);
