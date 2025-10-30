interface CategoryCardBackgroundProps {
  comingSoon?: boolean;
}

export const CategoryCardBackground = ({
  comingSoon = false,
}: CategoryCardBackgroundProps) => (
  <div className="absolute inset-0 opacity-10">
    <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent" />
    {comingSoon && (
      <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent animate-shimmer" />
    )}
  </div>
);
