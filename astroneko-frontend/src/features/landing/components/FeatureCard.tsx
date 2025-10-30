interface FeatureCardProps {
  title: string;
  description: string;
  index: number;
}

export function FeatureCard({ title, description, index }: FeatureCardProps) {
  return (
    <div
      className="max-w-[720px] 2xl:max-w-[1250px] mx-auto p-2 py-4 px-4 md:p-6 2xl:py-14 2xl:px-12 rounded-xl 2xl:rounded-[30px] bg-black/40 backdrop-blur-xs border border-white/8 2xl:border-white/16 transition-all duration-300 opacity-0 translate-y-5 animate-fade-in"
      style={{
        animationDelay: `${0.6 + index * 0.2}s`,
      }}
    >
      <h3 className="text-xs sm:text-sm md:text-base lg:text-lg xl:text-xl font-semibold text-white mb-2 sm:mb-3 md:mb-4">
        {title}
      </h3>
      <h3
        className="text-[#F7C36D] text-sm sm:text-base md:text-lg lg:text-xl xl:text-2xl font-semibold mt-2 sm:mt-3 md:mt-4 animate-text-glow"
        style={{
          animationDelay: "0.5s",
        }}
      >
        {description}
      </h3>
    </div>
  );
}
