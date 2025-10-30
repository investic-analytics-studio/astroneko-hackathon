/**
 * Referral Activation Hook
 *
 * Handles referral code activation logic including:
 * - Form state management
 * - API calls for referral activation
 * - Error handling with user feedback
 * - Success callbacks
 */

import { useState, useCallback } from "react";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { activateReferral as activateReferralAPI } from "@/apis/auth";
import { useAuth } from "./useAuth";
import { useInvalidateReferralCodes } from "./useReferralCodes";
import { logger } from "@/lib/logger";
import { track } from "@/lib/amplitude";

// Constants
const ERROR_RESET_DELAY_MS = 5000;
const DEFAULT_ERROR_MESSAGE = "Invalid referral code";
const TOAST_STYLES = {
  className: "custom-success-toast",
  descriptionClassName: "text-[#A1A1AA]",
  actionButtonStyle: { color: "#000000" },
} as const;

// Type definitions
/**
 * Submission state for referral form
 */
export type ReferralSubmissionState =
  | "idle"
  | "submitting"
  | "success"
  | "error";

/**
 * Return type for useReferralActivation hook
 */
export interface UseReferralActivationReturn {
  referralCode: string;
  submissionState: ReferralSubmissionState;
  errorMessage: string;
  setReferralCode: (code: string) => void;
  handleSubmit: (e: React.FormEvent) => Promise<void>;
  reset: () => void;
}

// Helper functions
/**
 * Extracts error message from API error response
 * @param error - Error object from API
 * @returns User-friendly error message
 */
const extractErrorMessage = (error: unknown): string => {
  if (error && typeof error === "object" && "response" in error) {
    const axiosError = error as {
      response: { data: { status: { message: string | string[] } } };
    };
    if (axiosError.response?.data?.status?.message) {
      const messages = axiosError.response.data.status.message;
      return Array.isArray(messages) ? messages[0] : messages;
    }
  } else if (error instanceof Error) {
    return error.message;
  }
  return DEFAULT_ERROR_MESSAGE;
};

/**
 * Custom hook for managing referral code activation
 *
 * Handles the complete referral activation flow including form state,
 * API communication, error handling, and success callbacks.
 *
 * @param onSuccess - Optional callback to run after successful activation
 * @returns Referral activation state and handlers
 *
 * @example
 * ```tsx
 * function ReferralModal({ onClose }) {
 *   const {
 *     referralCode,
 *     submissionState,
 *     errorMessage,
 *     setReferralCode,
 *     handleSubmit,
 *     reset
 *   } = useReferralActivation(() => {
 *     onClose();
 *   });
 *
 *   return (
 *     <form onSubmit={handleSubmit}>
 *       <input
 *         value={referralCode}
 *         onChange={(e) => setReferralCode(e.target.value)}
 *       />
 *       {errorMessage && <p>{errorMessage}</p>}
 *       <button disabled={submissionState === "submitting"}>
 *         Submit
 *       </button>
 *     </form>
 *   );
 * }
 * ```
 */
export function useReferralActivation(
  onSuccess?: () => void
): UseReferralActivationReturn {
  const { t } = useTranslation();
  const [referralCode, setReferralCode] = useState("");
  const [submissionState, setSubmissionState] =
    useState<ReferralSubmissionState>("idle");
  const [errorMessage, setErrorMessage] = useState("");
  const { refreshUserInfo } = useAuth();
  const invalidateReferralCodes = useInvalidateReferralCodes();

  /**
   * Handles referral code submission
   *
   * @param e - Form event
   */
  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!referralCode.trim()) return;

      setSubmissionState("submitting");
      setErrorMessage("");

      try {
        // Call API to activate referral code
        await activateReferralAPI({ referral_code: referralCode });

        // Refresh user info to get updated referral status
        await refreshUserInfo();

        // Invalidate referral codes cache to refresh the data
        invalidateReferralCodes();

        // Track successful referral code activation
        track("referral code confirmed", {
          referral_code: referralCode,
          success: true,
          placement: "referral_modal",
        });

        toast.success(t("waitlist.welcome_youre_in"), TOAST_STYLES);

        setSubmissionState("success");
        reset();
        onSuccess?.();
      } catch (error: unknown) {
        logger.error("Failed to activate referral:", error);

        const errorMsg = extractErrorMessage(error);
        setErrorMessage(errorMsg);
        setSubmissionState("error");

        // Reset error state after delay
        setTimeout(() => {
          setSubmissionState("idle");
          setErrorMessage("");
        }, ERROR_RESET_DELAY_MS);
      }
    },
    [referralCode, refreshUserInfo, invalidateReferralCodes, onSuccess]
  );

  /**
   * Resets form state to initial values
   */
  const reset = useCallback(() => {
    setReferralCode("");
    setSubmissionState("idle");
    setErrorMessage("");
  }, []);

  return {
    referralCode,
    submissionState,
    errorMessage,
    setReferralCode,
    handleSubmit,
    reset,
  };
}
