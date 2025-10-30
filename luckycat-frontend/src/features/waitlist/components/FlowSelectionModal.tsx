import React, { useCallback } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogTrigger,
} from "@/components/ui/dialog";
import { FlowSelection } from "./FlowSelection";
import { Button } from "@/components/ui/button";
import { useUserLimit, useAuth } from "@/hooks";
import { useTranslation } from "react-i18next";

interface FlowSelectionModalProps {
  trigger?: React.ReactNode;
  isOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
  onSelectReferral?: () => void;
  onSelectWaitlist?: () => void;
  onGoHome?: () => void;
}

export function FlowSelectionModal({
  trigger,
  isOpen,
  onOpenChange,
  onSelectReferral,
  onSelectWaitlist,
  onGoHome,
}: FlowSelectionModalProps) {
  const { authUser } = useAuth();
  const { isOverLimit } = useUserLimit();
  const { t } = useTranslation();

  // Check if user should be forced to waitlist
  const shouldForceWaitlist =
    authUser && !authUser.is_activated_referral && isOverLimit === true;

  const handleSelectReferral = () => {
    onOpenChange?.(false);
    onSelectReferral?.();
  };

  const handleSelectWaitlist = useCallback(() => {
    onOpenChange?.(false);
    onSelectWaitlist?.();
  }, [onOpenChange, onSelectWaitlist]);

  // No auto-redirect - let user see the limit message first

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      {trigger && <DialogTrigger asChild>{trigger}</DialogTrigger>}

      <DialogContent
        className="sm:max-w-md bg-black/95 border-white/20 rounded-[20px]"
        showCloseButton={false}
        onEscapeKeyDown={(e) => e.preventDefault()}
        onPointerDownOutside={(e) => e.preventDefault()}
      >
        <DialogHeader className="text-center">
          <DialogTitle className="text-2xl text-center font-bold text-[#F7C36D]">
            {t("waitlist.get_access")}
          </DialogTitle>
          <DialogDescription className="text-[#A1A1AA] text-center">
            {shouldForceWaitlist
              ? t("waitlist.quota_reached")
              : t("waitlist.choose_access")}
          </DialogDescription>
        </DialogHeader>

        {shouldForceWaitlist ? (
          <div className="space-y-6 mt-2">
            {/* Disabled Referral Option */}
            <Button
              disabled
              className="w-full bg-white/10 text-white/50 font-semibold border-white/10 rounded-[12px]
                       py-4 h-auto cursor-not-allowed opacity-50"
            >
              <div className="flex flex-col items-center gap-1">
                <span className="text-[16px] font-bold">
                  {t("waitlist.have_referral_code")}
                </span>
                <span className="text-[14px] font-medium">
                  {t("waitlist.currently_unavailable")}
                </span>
              </div>
            </Button>

            {/* Available Waitlist Option */}
            <Button
              onClick={handleSelectWaitlist}
              className="w-full bg-[#F7C36D] text-black font-semibold border-none rounded-[12px]
                       py-4 h-auto hover:bg-[#FFB53A] hover:border-none focus:border-none focus:ring-0 focus:outline-none transition-all duration-300"
            >
              <div className="flex flex-col items-center gap-1">
                <span className="text-[16px] font-bold">
                  {t("waitlist.join_waitlist")}
                </span>
                <span className="text-[14px] font-medium">
                  {t("waitlist.notify_ready")}
                </span>
              </div>
            </Button>
          </div>
        ) : (
          <FlowSelection
            onSelectReferral={handleSelectReferral}
            onSelectWaitlist={handleSelectWaitlist}
          />
        )}
        <div className="space-y-3">
          <Button
            onClick={onGoHome}
            className="w-full h-[46px] bg-transparent hover:bg-white/20 border-none rounded-md
                     text-white/70 hover:text-white font-medium transition-all duration-300
                     focus:border-none focus:ring-0 focus:outline-none"
          >
            {t("waitlist.go_to_home")}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

export default FlowSelectionModal;
