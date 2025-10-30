export const TypingIndicator = () => (
  <div className="flex gap-1.5">
    {[1, 2, 3].map((dot) => (
      <div
        key={dot}
        className="w-2 h-2 rounded-full bg-white animate-bounce-dot"
        style={{
          animationDelay: `${dot * 0.4}s`,
        }}
      />
    ))}
  </div>
);

export const TextOnlyLoadingSkeleton = () => (
  <div className="flex justify-start mb-12 2xl:mb-[200px] md:mx-20 opacity-0 animate-fade-in">
    <div className="relative max-w-2xl px-4 py-3 rounded-xl">
      <div className="p-3 rounded-lg border-none">
        <TypingIndicator />
      </div>
    </div>
  </div>
);
