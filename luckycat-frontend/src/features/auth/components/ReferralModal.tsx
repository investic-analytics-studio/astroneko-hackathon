import React from "react";
import { useTranslation } from "react-i18next";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogTrigger,
} from "@/components/ui/dialog";
import { ReferralForm } from "@/features/waitlist";
import { useReferralActivation } from "@/hooks";
import { AuthUser } from "@/apis/auth";
import { ReferralAlreadyActivated } from "./ReferralAlreadyActivated";

interface ReferralModalProps {
  trigger?: React.ReactNode;
  isOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
  onSuccess?: () => void;
  user?: AuthUser | null;
}

export function ReferralModal({
  trigger,
  isOpen,
  onOpenChange,
  onSuccess,
  user,
}: ReferralModalProps) {
  const { t } = useTranslation();

  // Use custom hook for referral activation logic
  const {
    referralCode,
    submissionState,
    errorMessage,
    setReferralCode,
    handleSubmit,
    reset,
  } = useReferralActivation(() => {
    onSuccess?.();
    onOpenChange?.(false);
  });

  // If user has already activated referral, show success message
  if (user?.is_activated_referral) {
    return (
      <Dialog open={isOpen} onOpenChange={onOpenChange}>
        {trigger && <DialogTrigger asChild>{trigger}</DialogTrigger>}

        <DialogContent
          className="sm:max-w-md bg-black/95 border-white/20 rounded-[20px]"
          showCloseButton={false}
        >
          <DialogHeader className="text-center">
            <DialogTitle className="text-2xl text-center font-bold text-[#F7C36D]">
              {t("referral_already_activated.title")}
            </DialogTitle>
          </DialogHeader>

          <ReferralAlreadyActivated onClose={() => onOpenChange?.(false)} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog
      open={isOpen}
      onOpenChange={(open) => {
        if (!open) {
          reset();
        }
        onOpenChange?.(open);
      }}
    >
      {trigger && <DialogTrigger asChild>{trigger}</DialogTrigger>}

      <DialogContent
        className="sm:max-w-md bg-black/95 border-white/20 rounded-[20px]"
        showCloseButton={true}
      >
        <DialogHeader className="text-center">
          <DialogTitle className="text-2xl text-center font-bold text-[#F7C36D]">
            {t("waitlist.enter_referral_title")}
          </DialogTitle>
          <DialogDescription className="text-[#A1A1AA] text-center">
            {t("waitlist.enter_referral_description")}
          </DialogDescription>
        </DialogHeader>

        <ReferralForm
          referralCode={referralCode}
          onReferralCodeChange={setReferralCode}
          onSubmit={handleSubmit}
          submissionState={submissionState}
          errorMessage={errorMessage}
        />
      </DialogContent>
    </Dialog>
  );
}

export default ReferralModal;
