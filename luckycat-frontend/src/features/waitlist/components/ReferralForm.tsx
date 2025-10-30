import React from "react";
import { Button } from "@/components/ui/button";
import { DialogFooter } from "@/components/ui/dialog";
import { ReferralCodeInfo } from "./ReferralCodeInfo";
import { useUserLimit, useAuth } from "@/hooks";
import { useTranslation } from "react-i18next";

type SubmissionState = "idle" | "submitting" | "success" | "error";

interface ReferralFormProps {
  referralCode: string;
  onReferralCodeChange: (code: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  submissionState: SubmissionState;
  errorMessage?: string;
}

export function ReferralForm({
  referralCode,
  onReferralCodeChange,
  onSubmit,
  submissionState,
  errorMessage,
}: ReferralFormProps) {
  const { authUser } = useAuth();
  const { isOverLimit } = useUserLimit();
  const { t } = useTranslation();

  // Check if user should be blocked from using referral codes
  const isBlockedFromReferral =
    authUser && !authUser.is_activated_referral && isOverLimit === true;

  return (
    <form onSubmit={onSubmit} className="space-y-6 mt-2">
      <div className="space-y-3">
        <label
          htmlFor="referralCode"
          className="block text-sm font-medium text-white"
        >
          {t("waitlist.referral_code")}
        </label>
        <input
          type="text"
          id="referralCode"
          value={referralCode}
          onChange={(e) => onReferralCodeChange(e.target.value)}
          className={`w-full px-4 py-3 rounded-lg text-white
                   focus:outline-none focus:ring-0 focus:border-none
                   placeholder:text-white/25 transition-colors ${
                     isBlockedFromReferral
                       ? "bg-white/10 cursor-not-allowed opacity-50"
                       : "bg-white/15"
                   }`}
          placeholder={
            isBlockedFromReferral
              ? t("waitlist.referral_unavailable")
              : t("waitlist.enter_referral_code")
          }
          required
          disabled={submissionState === "submitting" || !!isBlockedFromReferral}
        />
        {isBlockedFromReferral && (
          <p className="text-yellow-500 text-sm">
            {t("waitlist.quota_message")}
          </p>
        )}
        {submissionState === "error" && errorMessage && (
          <p className="text-red-500 text-sm">{errorMessage}</p>
        )}
      </div>

      <DialogFooter className="flex flex-col gap-3">
        <Button
          type="submit"
          disabled={
            submissionState === "submitting" ||
            !referralCode.trim() ||
            !!isBlockedFromReferral
          }
          className="w-full h-[46px] bg-[#F7C36D] hover:bg-[#FFB53A] border-none rounded-md
          hover:border-none focus:border-none focus:ring-0 focus:outline-none
          transition-all duration-300 text-black font-semibold
          disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {submissionState === "submitting" ? (
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 border-2 border-black/20 border-t-black rounded-full animate-spin" />
              {t("waitlist.verifying")}
            </div>
          ) : (
            t("common.confirm")
          )}
        </Button>
      </DialogFooter>
      <ReferralCodeInfo />
    </form>
  );
}
