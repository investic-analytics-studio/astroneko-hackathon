import { useMutation } from "@tanstack/react-query";
import { activateReferral } from "@/apis/auth";
import { useAuthStore } from "@/store/authStore";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { logger } from "@/lib/logger";

/**
 * Hook for handling referral code activation
 * Focuses solely on the activation API call and state management
 */
export function useReferralActivation() {
  const { t } = useTranslation();
  const { refreshUserInfo } = useAuthStore();

  const mutation = useMutation({
    mutationFn: activateReferral,
    onSuccess: async () => {
      logger.info("Referral activation successful");

      // Refresh user info to get updated is_activated_referral status
      try {
        await refreshUserInfo();
        logger.info("User info refreshed after referral activation");
      } catch (error) {
        logger.error("Failed to refresh user info after activation:", error);
        // Continue with flow even if refresh fails, as activation was successful
      }

      toast.success(
        t("referral.activateSuccess", "Referral code activated successfully!")
      );
    },
    onError: (error: any) => {
      logger.error("Referral activation failed:", error);

      const errorMessage =
        error?.response?.data?.status?.message?.[0] ||
        error?.response?.data?.data?.message ||
        t("referral.activateError", "Failed to activate referral code");

      toast.error(errorMessage);
    },
  });

  return {
    activate: mutation.mutate,
    isLoading: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
  };
}
