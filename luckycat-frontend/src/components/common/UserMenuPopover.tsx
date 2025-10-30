import { Button } from "../ui/button";
import { UserAvatar } from "../ui/UserAvatar";
import { IconLogout } from "@tabler/icons-react";
import { AuthUser } from "../../apis/auth";
import { Binary, ChevronRight, KeyRound } from "lucide-react";
import { useTranslation } from "react-i18next";

interface UserMenuPopoverProps {
  user: AuthUser;
  onLogout: () => void;
  onCloseMenu?: () => void;
  onOpenMyReferral?: () => void;
  onOpenFillReferral?: () => void;
}

interface UserHeaderProps {
  user: AuthUser;
}

export const UserHeader = ({ user }: UserHeaderProps) => {
  return (
    <div className="flex items-center gap-3 pb-3 border-b border-white/20">
      <UserAvatar user={user} size="sm" showTooltip={false} />
      <div className="flex flex-col">
        <span className="font-semibold text-white">
          {user.display_name || user.email}
        </span>
        <span className="text-sm text-white/70">{user.email}</span>
      </div>
    </div>
  );
};

interface MenuItemProps {
  icon: React.ReactNode;
  label: string;
  onClick: () => void;
  onKeyDown?: (event: React.KeyboardEvent) => void;
}

export const MenuItem = ({
  icon,
  label,
  onClick,
  onKeyDown,
}: MenuItemProps) => {
  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      onClick();
    }
    onKeyDown?.(event);
  };

  return (
    <div
      role="button"
      tabIndex={0}
      className="flex items-center justify-between hover:bg-white/10 rounded-md cursor-pointer py-1 border-b border-white/20 focus:outline-none"
      onClick={onClick}
      onKeyDown={handleKeyDown}
      aria-label={label}
    >
      <div className="flex items-center gap-2 w-full text-white text-[16px] p-2">
        <div className="flex items-center justify-center bg-[#F7C36D]/20 rounded-sm w-6 h-6">
          {icon}
        </div>
        {label}
      </div>
      <div className="flex items-center justify-center w-6 h-6">
        <ChevronRight className="w-4 h-4 text-white/70" />
      </div>
    </div>
  );
};

interface MyReferralCodeProps {
  onOpenMyReferral: () => void;
  onCloseMenu: () => void;
}

export const MyReferralCode = ({
  onOpenMyReferral,
  onCloseMenu,
}: MyReferralCodeProps) => {
  const { t } = useTranslation();

  const handleClick = () => {
    onOpenMyReferral();
    onCloseMenu();
  };

  return (
    <MenuItem
      icon={<Binary className="w-4 h-4 text-[#F7C36D]" />}
      label={t("user_menu.my_referral_code")}
      onClick={handleClick}
    />
  );
};

interface FillReferralCodeProps {
  onOpenFillReferral: () => void;
  onCloseMenu: () => void;
}

export const FillReferralCode = ({
  onOpenFillReferral,
  onCloseMenu,
}: FillReferralCodeProps) => {
  const { t } = useTranslation();

  const handleClick = () => {
    onOpenFillReferral();
    onCloseMenu();
  };

  return (
    <MenuItem
      icon={<KeyRound className="w-4 h-4 text-[#F7C36D]" />}
      label={t("user_menu.fill_referral_code")}
      onClick={handleClick}
    />
  );
};

interface LogoutButtonProps {
  onLogout: () => void;
}

export const LogoutButton = ({ onLogout }: LogoutButtonProps) => {
  const { t } = useTranslation();

  return (
    <div>
      <Button
        onClick={onLogout}
        className="w-full h-[46px] xl:h-auto bg-red-600 rounded-full xl:rounded-md hover:bg-red-500 border-none text-white flex items-center mt-3 gap-2 justify-center
        focus:outline-none focus:border-none
        hover:border-none hover:outline-none
        transition-all duration-300"
        aria-label={t("user_menu.logout")}
      >
        <IconLogout className="w-4 h-4" />
        {t("user_menu.logout")}
      </Button>
    </div>
  );
};

export const UserMenuPopover = ({
  user,
  onLogout,
  onCloseMenu,
  onOpenMyReferral,
  onOpenFillReferral,
}: UserMenuPopoverProps) => {
  return (
    <div className="flex flex-col gap-1">
      <UserHeader user={user} />

      {user.is_activated_referral && onOpenMyReferral && onCloseMenu && (
        <MyReferralCode
          onOpenMyReferral={onOpenMyReferral}
          onCloseMenu={onCloseMenu}
        />
      )}

      {!user.is_activated_referral && onOpenFillReferral && onCloseMenu && (
        <FillReferralCode
          onOpenFillReferral={onOpenFillReferral}
          onCloseMenu={onCloseMenu}
        />
      )}

      <LogoutButton onLogout={onLogout} />
    </div>
  );
};
