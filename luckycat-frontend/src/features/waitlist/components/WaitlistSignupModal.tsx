import React, { useState, useEffect } from "react";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogTrigger,
} from "@/components/ui/dialog";
import { WaitlistForm } from "./WaitlistForm";
import { SuccessState } from "./SuccessState";
import { useAuth, useAppNavigation } from "@/hooks";
import { useWaitlistStore } from "../../../store/waitlistStore";
import { logger } from "@/lib/logger";
import { track } from "@/lib/amplitude";

interface WaitlistSignupModalProps {
  trigger?: React.ReactNode;
  isOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
  showSuccessState?: boolean;
}

type SubmissionState = "idle" | "checking" | "submitting" | "success" | "error";

export function WaitlistSignupModal({
  trigger,
  isOpen,
  onOpenChange,
  showSuccessState = false,
}: WaitlistSignupModalProps) {
  const { t } = useTranslation();
  const [submissionState, setSubmissionState] = useState<SubmissionState>(
    showSuccessState ? "success" : "idle"
  );
  const { authUser } = useAuth();
  const { joinWaitlist, checkWaitlistStatus } = useWaitlistStore();
  const { goToHome } = useAppNavigation();

  // Update submission state when showSuccessState prop changes
  useEffect(() => {
    setSubmissionState(showSuccessState ? "success" : "idle");
  }, [showSuccessState]);

  const handleWaitlistSubmit = async (email: string) => {
    setSubmissionState("checking");

    try {
      // First, check if user is already in waitlist
      const waitlistData = await checkWaitlistStatus(email);

      if (waitlistData?.is_in_waiting_list) {
        // User is already in waitlist, show success state
        setSubmissionState("success");

        track('waitlist already joined', {
          email_domain: email.split('@')[1]?.toLowerCase() || 'unknown',
          result: 'already_joined',
        });

        toast.info(t("waitlist.already_in_waitlist"), {
          className: "custom-info-toast",
          descriptionClassName: "text-[#A1A1AA]",
        });

        return;
      }

      // User is not in waitlist, proceed with joining
      setSubmissionState("submitting");
      await joinWaitlist(email);
      setSubmissionState("success");

      track('waitlist joined successfully', {
        email_domain: email.split('@')[1]?.toLowerCase() || 'unknown',
        result: 'success',
      });

      toast.success(t("waitlist.join_success"), {
        className: "custom-success-toast",
        descriptionClassName: "text-[#A1A1AA]",
        actionButtonStyle: {
          color: "#000000",
        },
      });
    } catch (error) {
      logger.error("Failed to process waitlist:", error);
      setSubmissionState("error");

      track('waitlist join failed', {
        error_message: error instanceof Error ? error.message : "Unknown error",
        error_type: error instanceof Error ? error.constructor.name : "Unknown",
      });

      const errorMessage =
        error instanceof Error ? error.message : "Failed to process waitlist";
      toast.error(errorMessage);

      // Reset error state after 3 seconds
      setTimeout(() => {
        setSubmissionState("idle");
      }, 3000);
    }
  };

  const handleGoHome = () => {
    onOpenChange?.(false);
    goToHome();
  };

  const handleAddAnotherEmail = () => {
    setSubmissionState("idle");
  };

  const reset = () => {
    setSubmissionState(showSuccessState ? "success" : "idle");
  };

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
            {submissionState === "success" ? "" : t("waitlist.join_waitlist")}
          </DialogTitle>
          <DialogDescription className="text-[#A1A1AA] text-center">
            {submissionState === "success" ? "" : t("waitlist.notify_ready")}
          </DialogDescription>
        </DialogHeader>

        {submissionState === "success" ? (
          <SuccessState
            onGoHome={handleGoHome}
            onAddAnotherEmail={handleAddAnotherEmail}
          />
        ) : (
          <WaitlistForm
            userEmail={authUser?.email || ""}
            onSubmit={handleWaitlistSubmit}
            submissionState={submissionState}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}

export default WaitlistSignupModal;
