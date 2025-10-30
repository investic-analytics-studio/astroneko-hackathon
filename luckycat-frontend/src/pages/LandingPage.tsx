import { useState } from "react";
import { useImagePreloader, useAuth } from "@/hooks";
import { useTranslation } from "react-i18next";
import {
  HeroSection,
  ActionButtons,
  FeatureCard,
  VersionInfo,
} from "@/features/landing";
import { ReferralModal, LoginModal } from "@/features/auth";

export default function LandingPage() {
  const { t } = useTranslation();
  const { authUser } = useAuth();
  const [isFillReferralModalOpen, setIsFillReferralModalOpen] = useState(false);
  const [isLoginModalOpen, setIsLoginModalOpen] = useState(false);

  const handleOpenReferralModal = () => {
    setIsFillReferralModalOpen(true);
  };

  const handleShowLogin = () => {
    setIsLoginModalOpen(true);
  };

  useImagePreloader({
    images: [
      "bg/bg-maneki-4.webp",
      "bg/bg-maneki-5.webp",
      "bg/bg-maneki-6.webp",
    ],
    priority: false,
  });

  return (
    <div className="min-h-screen relative bg-[image:var(--bg-display-2)] bg-cover bg-center overflow-hidden flex flex-col items-center justify-center opacity-0 animate-fade-in">
      <div className="container relative z-10 mt-14 sm:mt-0 px-4 py-16 space-y-1 text-center">
        <HeroSection />
        <ActionButtons
          onActionClick={handleOpenReferralModal}
          onShowLogin={handleShowLogin}
        />

        <div className="font-press-start grid grid-cols-1 md:grid-cols-1 gap-6 mx-auto mt-10 2xl:mt-20 2xl:px-20 opacity-0 animate-fade-in-delayed">
          {[
            {
              title: t("landing.feature_title"),
              description: t("landing.feature_description"),
            },
          ].map((feature, index) => (
            <FeatureCard
              key={index}
              title={feature.title}
              description={feature.description}
              index={index}
            />
          ))}
          <VersionInfo />
        </div>
      </div>

      <ReferralModal
        isOpen={isFillReferralModalOpen}
        onOpenChange={setIsFillReferralModalOpen}
        user={authUser}
      />

      {/* Login Modal */}
      <LoginModal
        isOpen={isLoginModalOpen}
        onOpenChange={setIsLoginModalOpen}
      />
    </div>
  );
}
