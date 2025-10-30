import React, { useState } from "react";
import { DialogFooter } from "@/components/ui/dialog";
import { LoadingButton } from "@/components/ui/loading-button";
import { useTranslation } from "react-i18next";
import { track } from "@/lib/amplitude";

type SubmissionState = "idle" | "checking" | "submitting" | "success" | "error";

interface WaitlistFormProps {
  userEmail: string;
  onSubmit: (email: string) => void;
  submissionState: SubmissionState;
}

export function WaitlistForm({
  userEmail,
  onSubmit,
  submissionState,
}: WaitlistFormProps) {
  const { t } = useTranslation();
  const [email, setEmail] = useState(userEmail || "");
  const [emailError, setEmailError] = useState("");

  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    // Reset error
    setEmailError("");

    // Validate email
    if (!email.trim()) {
      setEmailError(t("waitlist.email_required") || "Email is required");
      return;
    }

    if (!validateEmail(email)) {
      setEmailError(
        t("waitlist.email_invalid") || "Please enter a valid email address"
      );
      return;
    }

    // Track waitlist email submission
    track('waitlist email submitted', {
      email_domain: email.split('@')[1]?.toLowerCase() || 'unknown',
      submission_source: 'modal',
    });

    onSubmit(email);
  };

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
    if (emailError) {
      setEmailError("");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6 mt-2">
      <div className="space-y-3">
        <label htmlFor="email" className="block text-sm font-medium text-white">
          {t("waitlist.email_address")}
        </label>
        <input
          type="email"
          id="email"
          value={email}
          onChange={handleEmailChange}
          placeholder={
            t("waitlist.email_placeholder") || "Enter your email address"
          }
          className="w-full px-4 py-3 bg-white/10 border border-white/10 rounded-lg text-white
                   focus:outline-none focus:ring-2 focus:ring-[#F7C36D] focus:border-transparent
                   placeholder:text-white/25 transition-colors"
          disabled={submissionState === "submitting"}
        />
        {(emailError || submissionState === "error") && (
          <p className="text-red-500 text-sm">
            {emailError || t("waitlist.failed_join")}
          </p>
        )}
      </div>

      <DialogFooter className="flex flex-col sm:flex-row gap-3">
        <LoadingButton
          type="submit"
          loading={
            submissionState === "checking" || submissionState === "submitting"
          }
          loadingText={
            submissionState === "checking"
              ? t("waitlist.checking")
              : t("waitlist.joining")
          }
          spinnerSize="sm"
          className="w-full h-[46px] bg-[#F7C36D] hover:bg-[#FFB53A] border-none rounded-md
          hover:border-none focus:border-none focus:ring-0 focus:outline-none
          transition-all duration-300 text-black font-semibold"
        >
          {t("waitlist.join_waitlist_btn")}
        </LoadingButton>
      </DialogFooter>
    </form>
  );
}
