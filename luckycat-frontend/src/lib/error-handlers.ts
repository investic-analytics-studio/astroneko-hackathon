import { toast } from "sonner";
import type { FortuneError, LimitInfo } from "@/apis/fortune";

const TOAST_STYLES = {
  className: "custom-error-toast",
  descriptionClassName: "text-[#A1A1AA]",
  actionButtonStyle: { color: "#000000" },
} as const;

export const formatLimitMessage = (
  key: string,
  limitInfo: LimitInfo,
  t: (key: string) => string
): string => {
  return t(key)
    .replace(/\{\{limit\}\}/g, String(limitInfo.limit))
    .replace(/\{\{used\}\}/g, String(limitInfo.used))
    .replace(/\{\{hours\}\}/g, String(limitInfo.resetHours))
    .replace(/\{\{mins\}\}/g, String(limitInfo.resetMins));
};

const showErrorToast = (message: string) => {
  toast.error(message, TOAST_STYLES);
};

export const handleFortuneError = (
  error: unknown,
  t: (key: string) => string,
  defaultMessage: string = "An unexpected error occurred"
): string => {
  if (!error || typeof error !== "object" || !("type" in error)) {
    return defaultMessage;
  }

  const fortuneError = error as FortuneError;

  const handlers: Record<string, (e: FortuneError) => string> = {
    bad_request: (e) => {
      showErrorToast(e.message);
      return t("waitlist.check_message_try_again");
    },
    free_trial_limit: (e) => {
      showErrorToast(e.message);
      return e.message;
    },
    daily_limit: (e) => {
      const message = e.limitInfo
        ? formatLimitMessage(
            "waitlist.daily_limit_exceeded_with_time",
            e.limitInfo,
            t
          )
        : t("waitlist.daily_limit_exceeded");
      showErrorToast(message);
      return message;
    },
    daily_limit_logged_in: (e) => {
      const message = e.limitInfo
        ? formatLimitMessage(
            "waitlist.daily_limit_logged_in_with_time",
            e.limitInfo,
            t
          )
        : t("waitlist.daily_limit_logged_in");
      showErrorToast(message);
      return message;
    },
    rate_limit: () => {
      const message = t("waitlist.rate_limit_exceeded");
      showErrorToast(message);
      return message;
    },
    server_error: (e) => {
      showErrorToast(e.message);
      return t("waitlist.fortune_teller_issues");
    },
    service_unavailable: (e) => {
      showErrorToast(e.message);
      return t("waitlist.service_unavailable");
    },
    network: (e) => {
      showErrorToast(e.message);
      return t("waitlist.network_error");
    },
  };

  const handler = handlers[fortuneError.type];
  if (handler) {
    return handler(fortuneError);
  }

  showErrorToast(t("waitlist.something_wrong"));
  return defaultMessage;
};
