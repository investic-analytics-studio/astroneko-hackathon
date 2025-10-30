import { useTranslation } from "react-i18next";

export function ReferralCodeInfo() {
  const { t } = useTranslation();

  return (
    <div className="space-y-2 border-t border-white/10 pt-4">
      <h3 className="text-sm font-medium text-white/80">
        {t("waitlist.referral_info_title")}
      </h3>
      <div className="p-3 rounded-lg bg-white/5 border border-white/10">
        <div className="text-sm text-white/70 space-y-1">
          <p>• {t("waitlist.referral_info_friend")}</p>
          <p>• {t("waitlist.referral_info_admin")}</p>
        </div>
      </div>
    </div>
  );
}
