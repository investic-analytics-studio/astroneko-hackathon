interface LogoButtonProps {
  onScrollToHome?: () => void;
}

export const LogoButton = ({ onScrollToHome }: LogoButtonProps) => {
  return (
    <a
      href="#"
      className="flex-shrink-0 flex items-center justify-start gap-2 cursor-pointer transition-all duration-300 hover:scale-105 py-2 focus:outline-none rounded-md bg-transparent"
      onClick={(e) => {
        e.preventDefault();
        onScrollToHome?.();
      }}
      aria-label="Go to home"
      tabIndex={0}
    >
      <img
        src="/logo/astro-logo.webp"
        alt="logo"
        className="w-[100px] sm:w-[120px] md:w-[140px] lg:w-[160px] xl:w-[180px] 2xl:w-[200px] h-auto"
      />
    </a>
  );
};
