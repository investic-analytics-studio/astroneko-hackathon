import { useTranslation } from "react-i18next";

interface ReferralCodeDisplayProps {
  referralCode: string;
}

export function ReferralCodeDisplay({ referralCode }: ReferralCodeDisplayProps) {
  const { t } = useTranslation();

  if (!referralCode) return null;

  return (
    <div className="mb-6">
      <p className="text-white/80 text-sm mb-2">
        {t("referral.activatingCode", "Activating code:")}
      </p>
      <p
        className="text-[#F7C36D] font-mono text-lg md:text-xl bg-black/30 px-4 py-3 rounded-lg border border-white/10"
        style={{
          textShadow: "2px 2px 4px rgba(0, 0, 0, 0.4)",
        }}
      >
        {referralCode}
      </p>
    </div>
  );
}