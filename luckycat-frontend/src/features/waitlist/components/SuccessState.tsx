import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { Clock, Plus } from "lucide-react";

interface SuccessStateProps {
  onGoHome: () => void;
  onAddAnotherEmail?: () => void;
}

export function SuccessState({
  onGoHome,
  onAddAnotherEmail,
}: SuccessStateProps) {
  const { t } = useTranslation();

  return (
    <div className="text-center py-0 space-y-6">
      <div className="w-18 h-18 mx-auto flex items-center justify-center">
        <Clock className="w-full h-full text-[#E78562] animate-spin duration-2000" />
      </div>

      <div className="space-y-4">
        <div className="space-y-2">
          <h3 className="text-[18px] font-bold text-[#F7C36D]">
            {t("waitlist.success_title")}
          </h3>
          <p className="text-[16px] font-medium text-[#A1A1AA]">
            {t("waitlist.success_message")}
          </p>
        </div>

        <div className="space-y-3">
          <Button
            onClick={onGoHome}
            className="w-full h-[46px] bg-[#F7C36D] hover:bg-[#FFB53A] border-none rounded-md
                     text-black font-semibold transition-all duration-300
                     focus:border-none focus:ring-0 focus:outline-none"
          >
            {t("waitlist.go_to_home")}
          </Button>

          {onAddAnotherEmail && (
            <Button
              onClick={onAddAnotherEmail}
              variant="outline"
              className="w-full h-[46px] bg-transparent border border-white/20 rounded-md
                       text-white font-semibold transition-all duration-300
                       hover:bg-white/20 hover:border-white/30 hover:text-white
                       focus:border-none focus:ring-0 focus:outline-none"
            >
              <Plus className="w-4 h-4 mr-2" />
              {t("waitlist.add_another_email") || "Add another email"}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
