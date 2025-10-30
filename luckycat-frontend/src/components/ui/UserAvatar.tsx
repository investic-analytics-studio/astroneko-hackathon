import { useState } from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Popover, PopoverContent, PopoverTrigger } from "./popover";
import { Drawer, DrawerContent, DrawerTrigger } from "./drawer";
import { useDeviceDetection } from "@/hooks";
import { AuthUser } from "../../apis/auth";

interface UserAvatarProps {
  user?: AuthUser;
  name?: string;
  size?: "sm" | "md" | "lg" | "responsive";
  onClick?: () => void;
  showTooltip?: boolean;
  showPopover?: boolean;
  popoverContent?: (onClose: () => void) => React.ReactNode;
}

export const UserAvatar = ({
  user,
  name,
  size = "md",
  onClick,
  showTooltip = true,
  showPopover = false,
  popoverContent,
}: UserAvatarProps) => {
  // State must be at the top level, not inside conditional
  const [isOpen, setIsOpen] = useState(false);

  // Get display name and initial
  const displayName = user?.display_name || user?.email || name || "";
  const initial = displayName.charAt(0).toUpperCase();
  const profileImageUrl = user?.profile_image_url;
  const { isMobile, isTablet } = useDeviceDetection();

  // Size variants with responsive sizing
  const sizeClasses = {
    sm: "w-5 h-5 sm:w-6 sm:h-6 text-xs sm:text-sm",
    md: "w-6 h-6 sm:w-8 sm:h-8 text-sm sm:text-base",
    lg: "w-8 h-8 sm:w-12 sm:h-12 text-base sm:text-lg",
    responsive: "w-8 h-8 sm:w-10 sm:h-10 text-sm sm:text-base",
  };

  const avatarContent = (
    <div
      onClick={onClick}
      className={`
        ${sizeClasses[size]} 
        flex items-center justify-center 
        rounded-full 
        bg-[#E78562]
        text-white 
        font-semibold 
        cursor-pointer 
        transition-all 
        duration-300 
        border-none
        hover:opacity-80
        overflow-hidden
        focus:outline-none
      `}
    >
      {profileImageUrl ? (
        <img
          src={profileImageUrl}
          alt={displayName}
          className="w-full h-full object-cover"
        />
      ) : (
        initial
      )}
    </div>
  );

  const content = showTooltip ? (
    <Tooltip>
      <TooltipTrigger asChild>{avatarContent}</TooltipTrigger>
      <TooltipContent
        side="bottom"
        sideOffset={6}
        className="bg-white text-black"
      >
        <p>{displayName}</p>
      </TooltipContent>
    </Tooltip>
  ) : (
    avatarContent
  );

  const handleClose = () => {
    setIsOpen(false);
  };

  if (showPopover && popoverContent) {
    if (isMobile || isTablet) {
      return (
        <Drawer open={isOpen} onOpenChange={setIsOpen}>
          <DrawerTrigger asChild>{content}</DrawerTrigger>
          <DrawerContent className="bg-black border-t border-white/20">
            <div className="p-10">{popoverContent(handleClose)}</div>
          </DrawerContent>
        </Drawer>
      );
    }

    return (
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>{content}</PopoverTrigger>
        <PopoverContent
          side="bottom"
          align="end"
          sideOffset={6}
          className="w-auto bg-black border border-white/20 p-4 rounded-lg shadow-lg"
        >
          {popoverContent(handleClose)}
        </PopoverContent>
      </Popover>
    );
  }

  return content;
};
