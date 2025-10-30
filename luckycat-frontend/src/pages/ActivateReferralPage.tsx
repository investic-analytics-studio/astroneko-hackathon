import { useReferralActivationFlow } from "@/hooks/referral/useReferralActivationFlow";
import { ReferralLogo } from "@/components/referral/ReferralLogo";
import { ActivationCard } from "@/components/referral/ActivationCard";
import { LoginModal } from "@/features/auth/components/LoginModal";

export default function ActivateReferralPage() {
  const {
    loginModalOpen,
    referralCode,
    isLoading,
    isSuccess,
    isError,
    isNavigating,
    handleLoginModalClose,
    handleBackToHome,
  } = useReferralActivationFlow();

  return (
    <div
      className={`min-h-screen relative bg-[image:var(--bg-display-2)] bg-cover bg-center overflow-hidden flex flex-col items-center justify-center transition-opacity duration-300 ${
        isNavigating ? "opacity-0" : "opacity-100"
      }`}
    >
      <div className="container relative z-10 px-4 py-16 space-y-6 text-center">
        <ReferralLogo />
        <ActivationCard
          referralCode={referralCode}
          isLoading={isLoading}
          isSuccess={isSuccess}
          isError={isError}
          onBackToHome={handleBackToHome}
        />
      </div>

      <LoginModal
        isOpen={loginModalOpen}
        onOpenChange={handleLoginModalClose}
      />
    </div>
  );
}
