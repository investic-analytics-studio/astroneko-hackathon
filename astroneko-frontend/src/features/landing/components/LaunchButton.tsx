import { useAppNavigation } from "@/hooks";
import { track } from "@/lib/amplitude";
import { trackLaunchAppClick } from "@/config/ga-init";

export function LaunchButton() {
  const { goToCategory } = useAppNavigation();

  const handleLaunch = () => {
    track("launch app clicked", {
      placement: "hero",
      button_text: "Launch App",
    });

    trackLaunchAppClick("hero", "Launch App");

    goToCategory();
  };

  return (
    <div className="flex flex-row sm:flex-row items-center gap-4 font-press-start justify-center sm:justify-center pt-6 opacity-0 translate-y-5 animate-slide-in-delayed">
      <button
        onClick={handleLaunch}
        style={{
          boxShadow: "0 0 10px 0 rgba(0, 0, 0, 0.2)",
        }}
        className="group flex items-center justify-center font-semibold border-none w-[80%] md:w-[300px] lg:w-auto h-[50px] rounded-[40px] px-2 md:px-18 py-6 md:py-7 text-sm xl:text-[16px] gap-2 bg-[#E78562] text-[#000000] hover:opacity-80 hover:border-none focus:border-none focus:ring-0 focus:outline-none transition-all duration-300 animate-pulse-glow"
      >
        Launch App
      </button>
    </div>
  );
}
