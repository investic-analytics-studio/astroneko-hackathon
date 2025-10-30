export function ReferralLogo() {
  return (
    <div className="opacity-0 translate-y-5 animate-fade-in">
      <div className="flex items-center justify-center gap-4 mb-6">
        <div className="animate-wiggle"></div>
        <div className="flex items-end gap-1 py-4 md:py-0">
          <img
            src="/logo/astro-logo.webp"
            alt="logo"
            className="w-[280px] sm:w-[320px] md:w-[400px] lg:w-[500px] h-auto"
          />
        </div>
      </div>
    </div>
  );
}