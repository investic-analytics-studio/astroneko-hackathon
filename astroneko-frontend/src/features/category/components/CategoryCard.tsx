import { LucideIcon } from "lucide-react";
import { CategoryCardBackground } from "./CategoryCardBackground";
import { CategoryCardContent } from "./CategoryCardContent";
import { ComingSoonCardContent } from "./ComingSoonCardContent";
import { track } from "@/lib/amplitude";
import { trackCategorySelected } from "@/config/ga-init";

interface CategoryCardProps {
  id: string;
  subtitle: string;
  title: string;
  icon: LucideIcon;
  color: string;
  hoverColor: string;
  comingSoon?: boolean;
  index: number;
  onSelect: (categoryId: string) => void;
}

export function CategoryCard({
  id,
  subtitle,
  title,
  icon,
  color,
  hoverColor,
  comingSoon = false,
  index,
  onSelect,
}: CategoryCardProps) {
  if (comingSoon) {
    return (
      <div
        className="
          bg-white/50 relative cursor-not-allowed rounded-3xl px-6 py-3 md:p-6 lg:p-6 xl:p-8
          text-white shadow-xl transform transition-all duration-300
          hover:scale-100 hover:shadow-2xl group overflow-hidden
          min-h-auto md:min-h-[240px] lg:min-h-[320px] xl:min-h-[400px]
          flex flex-col justify-between opacity-90
          translate-y-5 animate-slide-in-delayed
        "
        style={{
          animationDelay: `${0.3 + index * 0.1}s`,
        }}
      >
        <ComingSoonCardContent title={title} subtitle={subtitle} icon={icon} />
        <div className="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
      </div>
    );
  }

  const handleSelect = () => {
    // Use English names for consistent tracking
    const englishNames: Record<string, string> = {
      general: 'General',
      crypto: 'Crypto',
      lover: 'Lover',
      tarot: 'Tarot'
    };

    const englishTitle = englishNames[id] || title;

    track(`category selected - ${englishTitle.toLowerCase()}`, {
      category_id: id,
      category_name: englishTitle,
      category_type: id.toLowerCase(), // general, lover, crypto
      category_subtitle: subtitle,
    });

    trackCategorySelected(id, subtitle);
    onSelect(id);
  };

  return (
    <div
      onClick={handleSelect}
      className={`
        ${color} ${hoverColor}
        relative cursor-pointer rounded-3xl px-6 py-3 md:p-6 lg:p-6 xl:p-8
        text-white shadow-xl transform transition-all duration-300
        hover:scale-105 hover:shadow-2xl group overflow-hidden
        min-h-auto md:min-h-[240px] lg:min-h-[320px] xl:min-h-[400px]
        flex flex-col justify-between
        opacity-0 translate-y-5 animate-slide-in-delayed
      `}
      style={{
        animationDelay: `${0.3 + index * 0.1}s`,
      }}
    >
      <CategoryCardBackground comingSoon={comingSoon} />
      <CategoryCardContent title={title} subtitle={subtitle} icon={icon} />
      <div className="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
    </div>
  );
}
