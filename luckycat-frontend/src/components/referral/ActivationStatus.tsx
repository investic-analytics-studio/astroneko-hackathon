import { Loader2, CheckCircle, XCircle } from "lucide-react";
import { useAuthStore } from "@/store/authStore";
import { useTranslation } from "react-i18next";
import { useEffect, useState } from "react";

interface ActivationStatusProps {
  isLoading: boolean;
  isSuccess: boolean;
  isError: boolean;
}

export function ActivationStatus({
  isLoading,
  isSuccess,
  isError,
}: ActivationStatusProps) {
  const { t } = useTranslation();
  const { authUser: user } = useAuthStore();
  const [showStatus, setShowStatus] = useState(false);

  useEffect(() => {
    // Reset and show status on props change for smooth transitions
    setShowStatus(false);
    const timer = setTimeout(() => setShowStatus(true), 100);
    return () => clearTimeout(timer);
  }, [isLoading, isSuccess, isError]);

  const getStatusContent = () => {
    if (isLoading) {
      return {
        icon: <Loader2 className="h-6 w-6 sm:h-8 sm:w-8" />,
        text: t("referral.activating", "Activating..."),
        className: "text-[#F7C36D]",
        bgColor: "bg-[#F7C36D]/10",
        borderColor: "border-[#F7C36D]/20",
        ringColor: "ring-[#F7C36D]/10",
      };
    }

    if (isSuccess) {
      return {
        icon: <CheckCircle className="h-8 w-8 sm:h-10 sm:w-10" />,
        text: t("referral.success", "Referral code activated!"),
        className: "text-[#10B981]",
        bgColor: "bg-[#10B981]/10",
        borderColor: "border-[#10B981]/20",
        ringColor: "ring-[#10B981]/10",
      };
    }

    if (isError) {
      return {
        icon: <XCircle className="h-8 w-8 sm:h-10 sm:w-10" />,
        text: t("referral.failed", "Activation failed"),
        className: "text-[#BD042D]",
        bgColor: "bg-[#BD042D]/10",
        borderColor: "border-[#BD042D]/20",
        ringColor: "ring-[#BD042D]/10",
      };
    }

    if (!user) {
      return {
        icon: <Loader2 className="h-6 w-6 sm:h-8 sm:w-8" />,
        text: t("referral.waitingForLogin", "Waiting for login..."),
        className: "text-[#E78562]",
        bgColor: "bg-[#E78562]/10",
        borderColor: "border-[#E78562]/20",
        ringColor: "ring-[#E78562]/10",
      };
    }

    return {
      icon: <Loader2 className="h-6 w-6 sm:h-8 sm:w-8" />,
      text: t("referral.processing", "Processing..."),
      className: "text-[#E78562]",
      bgColor: "bg-[#E78562]/10",
      borderColor: "border-[#E78562]/20",
      ringColor: "ring-[#E78562]/10",
    };
  };

  const status = getStatusContent();

  if (!user && !isLoading) {
    return (
      <div
        className={`flex items-center gap-2 transition-all duration-300 ${
          showStatus
            ? "opacity-100 translate-y-0"
            : "opacity-0 translate-y-2"
        }`}
      >
        <div
          className={`${isLoading ? "animate-spin" : ""} ${
            status.className
          } flex-shrink-0`}
        >
          {status.icon}
        </div>
        <span
          className={`text-base font-medium ${status.className}`}
        >
          {status.text}
        </span>
      </div>
    );
  }

  return (
    <div
      className={`flex items-center gap-2 transition-all duration-300 ${
        showStatus
          ? "opacity-100 translate-y-0"
          : "opacity-0 translate-y-2"
      }`}
    >
      <div
        className={`${
          isSuccess
            ? "animate-bounce"
            : isError
            ? "animate-pulse"
            : isLoading
            ? "animate-spin"
            : ""
        } ${status.className} flex-shrink-0`}
      >
        {status.icon}
      </div>
      <span
        className={`text-base font-medium ${status.className}`}
      >
        {status.text}
      </span>
    </div>
  );
}
