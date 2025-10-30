import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useState } from "react";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/hooks";
import { logger } from "@/lib/logger";
import { trackLoginClick } from "@/config/ga-init";

const GoogleIcon = () => {
  return (
    <svg
      className="w-5 h-5"
      viewBox="0 0 24 24"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
        fill="#4285F4"
      />
      <path
        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
        fill="#34A853"
      />
      <path
        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
        fill="#FBBC05"
      />
      <path
        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
        fill="#EA4335"
      />
    </svg>
  );
};

interface LoginModalProps {
  children?: React.ReactNode;
  isOpen?: boolean;
  onOpenChange?: (open: boolean) => void;
}

export const LoginModal = ({ children, isOpen: controlledIsOpen, onOpenChange }: LoginModalProps) => {
  const [internalIsOpen, setInternalIsOpen] = useState(false);
  const isControlled = controlledIsOpen !== undefined;
  const isOpen = isControlled ? controlledIsOpen : internalIsOpen;
  const setIsOpen = isControlled ? onOpenChange! : setInternalIsOpen;
  const { t } = useTranslation();
  const { loginWithGoogle } = useAuth();

  const handleGoogleLogin = async () => {
    try {
      const result = await loginWithGoogle();

      // If login was cancelled by user, don't show success or error
      if (!result) {
        return;
      }

      toast.success(t("waitlist.login_successful"), {
        className: "custom-success-toast",
        descriptionClassName: "text-[#A1A1AA]",
        actionButtonStyle: {
          color: "#000000",
        },
      });

      trackLoginClick("hero", "Login Click");
      setIsOpen(false);
    } catch (error) {
      logger.error("Login failed:", error);

      const errorMessage =
        error instanceof Error ? error.message : t("waitlist.login_failed");
      toast.error(errorMessage);
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      {children && <DialogTrigger asChild>{children}</DialogTrigger>}
      <DialogContent className="sm:max-w-[460px] bg-black/95 backdrop-blur-xl border-white/20 rounded-[20px] text-white p-10">
        <DialogHeader>
          <DialogTitle className="text-2xl font-bold text-center text-white">
            Sign in
          </DialogTitle>
          <DialogDescription className="text-center text-gray-300">
            Sign in to your account to continue your fortune journey
          </DialogDescription>
        </DialogHeader>

        <div className="mt-2">
          {/* Google Sign-In Button */}
          <div className="w-full">
            <Button
              type="button"
              onClick={handleGoogleLogin}
              className="w-full h-[46px] bg-white hover:bg-white/90 hover:border-none  active:scale-98 text-black text-[14px] font-medium py-3 px-4 rounded-md transition-all duration-300 flex items-center justify-center gap-3"
            >
              <>
                <GoogleIcon />
                Continue with Google
              </>
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
