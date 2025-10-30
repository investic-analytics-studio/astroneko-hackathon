import { useTranslation } from "react-i18next";
import { IconBrandX } from "@tabler/icons-react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

export const SocialLinks = () => {
  const { t } = useTranslation();

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <a
          href="https://x.com/masterastroneko"
          target="_blank"
          rel="noopener noreferrer"
          className="w-8 h-8 sm:w-9 sm:h-9 md:w-10 md:h-10 lg:w-11 lg:h-11 flex items-center justify-center font-medium text-black bg-white rounded-lg md:rounded-xl hover:bg-white/90 hover:scale-105 active:scale-95 transition-all duration-300 cursor-pointer focus:outline-none focus:ring-2 focus:ring-white/50"
          aria-label={t("common.follow_us_on_x")}
        >
          <IconBrandX className="w-3 h-3 sm:w-4 sm:h-4 md:w-5 md:h-5" />
        </a>
      </TooltipTrigger>
      <TooltipContent
        side="bottom"
        sideOffset={6}
        className="bg-white text-black text-xs sm:text-sm md:text-base"
      >
        <p>{t("common.follow_us_on_x")}</p>
      </TooltipContent>
    </Tooltip>
  );
};
