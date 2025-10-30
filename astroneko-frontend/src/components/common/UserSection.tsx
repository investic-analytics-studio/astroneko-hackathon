import { LoginModal, MyReferralModal, ReferralModal } from "@/features/auth";
import { User } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { AuthUser } from "../../apis/auth";
import { Button } from "../ui/button";
import { UserAvatar } from "../ui/UserAvatar";
import { UserMenuPopover } from "./UserMenuPopover";

interface UserMenuPopoverContentProps {
  user: AuthUser;
  onLogout: () => void;
  onClose: () => void;
  onOpenMyReferral: () => void;
  onOpenFillReferral: () => void;
}

export const UserMenuPopoverContent = ({
  user,
  onLogout,
  onClose,
  onOpenMyReferral,
  onOpenFillReferral,
}: UserMenuPopoverContentProps) => {
  return (
    <UserMenuPopover
      user={user}
      onLogout={onLogout}
      onCloseMenu={onClose}
      onOpenMyReferral={onOpenMyReferral}
      onOpenFillReferral={onOpenFillReferral}
    />
  );
};

interface UserSectionProps {
  authUser: AuthUser | null;
  isLoading: boolean;
  onLogout: () => void;
}

export const UserSection = ({
  authUser,
  isLoading,
  onLogout,
}: UserSectionProps) => {
  const { t } = useTranslation();

  const [isMyReferralModalOpen, setIsMyReferralModalOpen] = useState(false);
  const [isFillReferralModalOpen, setIsFillReferralModalOpen] = useState(false);

  const renderUserContent = () => {
    if (isLoading) {
      return (
        <div className="flex items-center gap-1 sm:gap-2">
          <div className="w-8 h-8 sm:w-9 sm:h-9 md:w-10 md:h-10 lg:w-11 lg:h-11 bg-white/20 rounded-full animate-pulse" />
        </div>
      );
    }

    if (authUser) {
      return (
        <div className="flex items-center gap-1 sm:gap-2">
          <UserAvatar
            user={authUser}
            size="responsive"
            showTooltip={false}
            showPopover={true}
            popoverContent={(onClose) => (
              <UserMenuPopoverContent
                user={authUser}
                onLogout={onLogout}
                onClose={onClose}
                onOpenMyReferral={() => setIsMyReferralModalOpen(true)}
                onOpenFillReferral={() => setIsFillReferralModalOpen(true)}
              />
            )}
          />
        </div>
      );
    }

    return (
      <LoginModal>
        <div>
          <Button
            className="px-2 sm:px-3 h-8 sm:h-9 md:h-10 lg:h-11 text-xs sm:text-sm font-semibold text-black bg-[#F7C36D] hover:bg-[#F7C36D]/80 rounded-lg md:rounded-xl hover:scale-105 transition-all duration-200
          focus:none ring-0 outline-none border-0"
          >
            <User className="w-4 h-4 sm:hidden" />
            <span className="hidden sm:inline px-2">{t("common.login")}</span>
          </Button>
        </div>
      </LoginModal>
    );
  };

  return (
    <>
      {renderUserContent()}

      <MyReferralModal
        isOpen={isMyReferralModalOpen}
        onOpenChange={setIsMyReferralModalOpen}
      />
      <ReferralModal
        isOpen={isFillReferralModalOpen}
        onOpenChange={setIsFillReferralModalOpen}
      />
    </>
  );
};
