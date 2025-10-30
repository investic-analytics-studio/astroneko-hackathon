import { LucideIcon, ChevronRight } from "lucide-react";

interface CategoryCardContentProps {
  title: string;
  subtitle: string;
  icon: LucideIcon;
}

export const CategoryCardContent = ({
  title,
  subtitle,
  icon: IconComponent,
}: CategoryCardContentProps) => (
  <div className="relative z-10 flex flex-row md:flex-col justify-between h-full gap-4 md:gap-6 xl:gap-16">
    <div className="flex flex-col md:flex-row gap-2 sm:gap-4 md:gap-6 lg:gap-8 xl:gap-10 items-start justify-between">
      <div className="flex flex-col items-start justify-start">
        <p className="text-xs sm:text-sm md:text-base lg:text-lg font-medium opacity-90 leading-relaxed">
          {subtitle}
        </p>
        <h3 className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold font-press-start">
          {title}
        </h3>
      </div>
      <div className="bg-white/20 rounded-full p-2 group-hover:bg-white/30 group-hover:rotate-90 transition-all duration-300">
        <ChevronRight className="w-5 h-5 md:w-6 md:h-6" />
      </div>
    </div>

    <div className="flex items-center justify-center">
      <div className="bg-white/20 rounded-full p-4 sm:p-4 md:p-6 lg:p-8 xl:p-10 transition-all ease-in-out duration-300 group-hover:bg-white/30 group-hover:scale-105 group-hover:rotate-5">
        <IconComponent className="w-6 h-6 md:w-16 md:h-16 lg:w-20 lg:h-20 xl:w-20 xl:h-20" />
      </div>
    </div>
  </div>
);
