import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useLanguage } from "@/hooks";
import { CategoryTitle, CategoryCard } from "@/features/category";
import { getCategories } from "../constants/categories";

export default function CategoryPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const { currentLanguage } = useLanguage();

  const handleCategorySelect = (categoryId: string) => {
    const categories = getCategories(t);
    const category = categories.find((cat) => cat.id === categoryId);

    // Don't navigate if it's coming soon
    if (category?.comingSoon) {
      return;
    }

    // Navigate to category-specific chat route
    navigate({
      to: "/$lng/chat/$category",
      params: {
        lng: currentLanguage,
        category: categoryId,
      },
    });
  };

  return (
    <div className="min-h-screen relative bg-[image:var(--bg-display-2)] bg-cover bg-center flex flex-col items-center opacity-0 animate-fade-in pt-20 md:pt-24 pb-8">
      <div className="container relative z-10 px-4 py-8 space-y-1 text-center">
        {/* Title with continuous fade-in animation */}
        <div className="opacity-0 animate-slide-in">
          <CategoryTitle />
        </div>

        {/* Category Cards with continuous staggered animations */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-4 gap-3 sm:gap-6 md:gap-6 lg:gap-6 xl:gap-8 md:max-w-2xl lg:max-w-5xl xl:max-w-7xl mx-auto mt-8">
          {getCategories(t).map((category, index) => (
            <div
              key={category.id}
              className="opacity-0 animate-fade-in-stagger"
              style={{
                animationDelay: `${index * 150 + 300}ms`,
                animationFillMode: "forwards",
              }}
            >
              <CategoryCard
                id={category.id}
                subtitle={category.subtitle}
                title={category.title}
                icon={category.icon}
                color={category.color}
                hoverColor={category.hoverColor}
                comingSoon={category.comingSoon}
                index={index}
                onSelect={handleCategorySelect}
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
