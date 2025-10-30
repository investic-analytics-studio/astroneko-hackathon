import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";

interface ReferralAlreadyActivatedProps {
  onClose: () => void;
}

export function ReferralAlreadyActivated({
  onClose,
}: ReferralAlreadyActivatedProps) {
  const { t } = useTranslation();

  return (
    <div className="text-center py-6">
      <div className="text-green-400 text-6xl mb-4">✓</div>
      <p className="text-white text-lg mb-8">
        {t("referral_already_activated.success_message")}
      </p>
      <Button
        onClick={onClose}
        className="hover:bg-[#E78562]/80 w-full font-semibold border-none h-[50px] rounded-[40px] px-4 md:px-6 py-6 md:py-7 text-sm xl:text-[16px] gap-2 bg-[#E78562] text-black hover:opacity-80 focus:outline-none transition-all duration-300 animate-pulse-glow"
      >
        {t("referral_already_activated.close_button")}
      </Button>
    </div>
  );
}
