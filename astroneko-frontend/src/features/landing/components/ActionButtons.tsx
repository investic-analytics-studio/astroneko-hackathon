import { useAppNavigation, useAuth } from "@/hooks";
import { Button } from "@/components/ui/button";
import { track } from "@/lib/amplitude";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import {
  trackLaunchAppClick,
  trackReferralModalClick,
  trackYourReferralCodeClick,
} from "@/config/ga-init";
import { MyReferralModal } from "@/features/auth/components/MyReferralModal";
import { useState } from "react";

interface ActionButtonsProps {
  onActionClick: () => void;
  onShowLogin: () => void;
}

export function ActionButtons({
  onActionClick,
  onShowLogin,
}: ActionButtonsProps) {
  const { goToCategory } = useAppNavigation();
  const { isAuthenticated, authUser } = useAuth();
  const { t } = useTranslation();
  const [isReferralModalOpen, setIsReferralModalOpen] = useState(false);

  const handleLaunchApp = () => {
    track("launch app clicked", {
      placement: "hero",
      button_text: "Launch App",
    });

    trackLaunchAppClick("hero", "Launch App");
    goToCategory();
  };

  const handleActionClick = () => {
    // If user has activated referral, track "your referral code clicked" events
    if (authUser?.is_activated_referral) {
      const buttonText = "your referral code";

      track("your referral code clicked", {
        placement: "hero",
        button_text: buttonText,
      });

      trackYourReferralCodeClick("hero", buttonText);
    } else {
      // Track regular referral modal events
      track("referral modal clicked", {
        placement: "hero",
        button_text: "referral modal",
      });

      trackReferralModalClick("hero", "referral modal");
    }

    // If user is not authenticated, show login modal instead
    if (!isAuthenticated) {
      // Show toast message informing user they need to login
      toast.info(t("waitlist.please_login"), {
        className: "custom-info-toast",
        descriptionClassName: "text-[#A1A1AA]",
        actionButtonStyle: {
          color: "#000000",
        },
      });

      onShowLogin();
      return;
    }

    // If user has activated referral, show their referral codes
    if (authUser?.is_activated_referral) {
      setIsReferralModalOpen(true);
      return;
    }

    onActionClick();
  };

  return (
    <div className="flex flex-col sm:flex-row items-center gap-4 lg:gap-6 font-press-start justify-center mt-6 md:mt-10  px-4 opacity-0 translate-y-5 animate-slide-in-delayed">
      {/* Launch App Button */}
      <Button
        onClick={handleLaunchApp}
        className="group flex items-center justify-center font-semibold border-none w-full sm:w-[200px] md:w-[220px] lg:w-[240px] xl:w-[260px] h-10 sm:h-[50px] rounded-[40px] px-4 md:px-6 py-6 md:py-7 text-sm xl:text-[16px] gap-2 bg-[#E78562] text-black hover:opacity-80 focus:outline-none transition-all duration-300 animate-pulse-glow"
        style={{
          boxShadow: "0 0 10px 0 rgba(0, 0, 0, 0.2)",
        }}
      >
        {t("action_buttons.launch_app")}
      </Button>

      {/* Referral Button */}
      <Button
        onClick={handleActionClick}
        className="group flex items-center justify-center font-semibold w-full sm:w-[200px] md:w-[220px] lg:w-[240px] xl:w-[260px] h-10 sm:h-[50px] rounded-[40px] px-4 md:px-6 py-6 md:py-7 text-sm xl:text-[16px] gap-2 bg-white/10 border-2 border-white/30 text-white hover:bg-white/20 focus:outline-none transition-all duration-300"
      >
        {authUser?.is_activated_referral
          ? t("action_buttons.invite_friends")
          : t("action_buttons.fill_code")}
      </Button>

      {/* My Referral Modal */}
      <MyReferralModal
        isOpen={isReferralModalOpen}
        onOpenChange={setIsReferralModalOpen}
      />
    </div>
  );
}
