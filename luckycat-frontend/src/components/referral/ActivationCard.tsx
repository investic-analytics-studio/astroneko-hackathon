import { UserCheck } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTranslation } from "react-i18next";
import { ReferralCodeDisplay } from "./ReferralCodeDisplay";
import { ActivationStatus } from "./ActivationStatus";

interface ActivationCardProps {
  referralCode: string;
  isLoading: boolean;
  isSuccess: boolean;
  isError: boolean;
  onBackToHome: () => void;
}

export function ActivationCard({
  referralCode,
  isLoading,
  isSuccess,
  isError,
  onBackToHome,
}: ActivationCardProps) {
  const { t } = useTranslation();

  return (
    <div
      className="max-w-[720px] mx-auto p-6 md:p-8 rounded-xl 2xl:rounded-[30px] bg-black/40 backdrop-blur-xs border border-white/16 transition-all duration-500 ease-out opacity-0 translate-y-5 animate-fade-in"
      style={{
        animationDelay: "0.1s",
      }}
    >
      <h2
        className="text-2xl md:text-3xl font-semibold text-white mb-6 flex items-center justify-center gap-3"
        style={{
          textShadow:
            "4px 4px 8px rgba(0, 0, 0, 0.4), 2px 2px 4px rgba(0, 0, 0, 0.4)",
        }}
      >
        <UserCheck className="h-8 w-8" />
        {t("referral.activateTitle", "Activate Referral Code")}
      </h2>

      <ReferralCodeDisplay referralCode={referralCode} />

      <div className="flex justify-center">
        <ActivationStatus
          isLoading={isLoading}
          isSuccess={isSuccess}
          isError={isError}
        />
      </div>

      <div className="text-center mt-6">
        <Button
          variant="ghost"
          onClick={onBackToHome}
          className="bg-[var(--brand-primary)] hover:bg-[var(--brand-primary)]/80 text-black font-semibold py-3 px-6 rounded-lg transition-all duration-300 focus:outline-none hover:outline-none"
        >
          {t("common.backToHome", "Back to Home")}
        </Button>
      </div>
    </div>
  );
}
