import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { Copy, Check, Share2, XIcon } from "lucide-react";
import { toast } from "sonner";
import { useState, useEffect, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { useReferralCodes } from "@/hooks";
import { useReferralStore } from "../../../store/referralStore";
import { ReferralCode } from "../../../apis/auth";
import { logger } from "@/lib/logger";

interface MyReferralModalProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

// Constants
const COPY_RESET_TIMEOUT = 5000;
const TOAST_CONFIG = {
  success: {
    className: "custom-success-toast",
    descriptionClassName: "text-[#A1A1AA]",
  },
  error: {
    className: "custom-error-toast",
    descriptionClassName: "text-[#A1A1AA]",
  },
} as const;

const DIALOG_CLASSES = {
  content: "sm:max-w-lg bg-black/95 border-white/20 rounded-[20px]",
  header: "text-center",
  title: "text-2xl text-center font-bold text-[#F7C36D]",
} as const;

const getReferralUrl = (code: string): string => {
  const baseUrl = window.location.origin;
  return `${baseUrl}/en/activate-referral-code?code=${code}`;
};

// Sub-components
interface ModalWrapperProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  children: React.ReactNode;
}

function ModalWrapper({
  isOpen,
  onOpenChange,
  title,
  description,
  children,
}: ModalWrapperProps) {
  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className={DIALOG_CLASSES.content} showCloseButton={false}>
        <DialogHeader className="relative space-y-2">
          <Button
            onClick={() => onOpenChange(false)}
            className="absolute -top-1 -right-1 rounded-full w-8 h-8 hover:bg-white/10 border-gray-700 z-10"
          >
            <XIcon className="w-4 h-4 text-white" />
          </Button>
          <DialogTitle className={DIALOG_CLASSES.title}>{title}</DialogTitle>
          {description && (
            <DialogDescription className="text-[#A1A1AA] text-center text-sm leading-relaxed">
              {description}
            </DialogDescription>
          )}
        </DialogHeader>
        {children}
      </DialogContent>
    </Dialog>
  );
}

function LoadingState() {
  return (
    <div className="flex items-center justify-center py-8">
      <Spinner size="lg" variant="accent" />
    </div>
  );
}

interface ErrorStateProps {
  message: string;
  errorDetails: string;
}

function ErrorState({ message, errorDetails }: ErrorStateProps) {
  return (
    <div className="text-center py-8">
      <p className="text-red-500">{message}</p>
      <p className="text-white/70 text-sm mt-2">{errorDetails}</p>
    </div>
  );
}

interface ShareReferralModalProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  referralCode: string;
}

function ShareReferralModal({
  isOpen,
  onOpenChange,
  referralCode,
}: ShareReferralModalProps) {
  const { t } = useTranslation();
  const [isCopied, setIsCopied] = useState(false);
  const referralUrl = getReferralUrl(referralCode);

  const handleCopyUrl = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(referralUrl);
      setIsCopied(true);

      setTimeout(() => setIsCopied(false), COPY_RESET_TIMEOUT);

      toast.success(t("waitlist.referral_url_copied"), TOAST_CONFIG.success);
    } catch (error) {
      logger.error("Failed to copy referral URL:", error);
      toast.error(t("waitlist.failed_copy_url"), TOAST_CONFIG.error);
    }
  }, [referralUrl, t]);

  return (
    <ModalWrapper
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      title={t("waitlist.share_referral_link")}
      description={t("waitlist.share_link_description")}
    >
      <div className="space-y-4 p-2">
        <div className="space-y-2">
          <h3 className="text-sm font-medium text-white/80">
            {t("waitlist.referral_link")}
          </h3>
          <div className="flex items-center gap-2 p-3 rounded-lg bg-white/5 border border-white/10">
            <div className="flex-1 overflow-hidden">
              <p className="text-sm text-white/90 break-all font-mono">
                {referralUrl}
              </p>
            </div>
            <Button
              onClick={handleCopyUrl}
              className="flex-shrink-0 h-9 px-3 bg-[#F7C36D] hover:bg-[#F7C36D]/90 text-black font-semibold rounded-lg flex items-center justify-center gap-2 transition-all duration-200"
            >
              {isCopied ? (
                <Check className="w-4 h-4" />
              ) : (
                <Copy className="w-4 h-4" />
              )}
            </Button>
          </div>
        </div>

        <div className="p-3 rounded-lg bg-white/5 border border-white/10">
          <p className="text-sm text-white/70">
            {t("waitlist.share_link_instruction")}
          </p>
        </div>
      </div>
    </ModalWrapper>
  );
}

interface ReferralCodeItemProps {
  referral: ReferralCode;
  index: number;
  isCopied: boolean;
  onCopy: (code: string, index: number) => void;
  onShare: (code: string) => void;
}

function ReferralCodeItem({
  referral,
  index,
  isCopied,
  onCopy,
  onShare,
}: ReferralCodeItemProps) {
  const { t } = useTranslation();

  return (
    <div className="flex items-center justify-between pl-3 sm:pl-4 pr-2 rounded-lg bg-white/5 border border-white/10 h-12 sm:h-14">
      <div className="flex items-center gap-2 sm:gap-4 flex-1 min-w-0">
        <div
          className={`flex items-center gap-1 sm:gap-2 text-[11px] sm:text-[14px] border-r border-white/15 pr-2 ${
            referral.is_activated ? "text-green-500" : "text-gray-500"
          }`}
        >
          <span className="font-medium whitespace-nowrap">
            {referral.is_activated
              ? t("waitlist.active")
              : t("waitlist.inactive")}
          </span>
        </div>
        <span className="font-mono text-base sm:text-2xl text-white font-bold tracking-wider truncate">
          {referral.referral_code}
        </span>
      </div>
      {!referral.is_activated && (
        <div className="flex items-center gap-1 flex-shrink-0">
          <Button
            onClick={() => onShare(referral.referral_code)}
            className="w-8 h-8 sm:w-9 sm:h-9 p-0 bg-transparent hover:bg-white/10 border-none text-white rounded-lg flex items-center justify-center transition-all duration-200 focus:outline-none focus:ring-0"
          >
            <Share2 className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
          </Button>
          <Button
            onClick={() => onCopy(referral.referral_code, index)}
            className="w-8 h-8 sm:w-9 sm:h-9 p-0 bg-transparent hover:bg-white/10 border-none text-white rounded-lg flex items-center justify-center transition-all duration-200 focus:outline-none focus:ring-0"
          >
            {isCopied ? (
              <Check className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-green-500" />
            ) : (
              <Copy className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
            )}
          </Button>
        </div>
      )}
    </div>
  );
}

interface ReferralCodesListProps {
  codes: ReferralCode[];
  copiedIndex: number | null;
  onCopy: (code: string, index: number) => void;
  onShare: (code: string) => void;
  emptyMessage: string;
  title: string;
}

function ReferralCodesList({
  codes,
  copiedIndex,
  onCopy,
  onShare,
  emptyMessage,
  title,
}: ReferralCodesListProps) {
  if (codes.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-white/70">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <h3 className="text-sm font-medium text-white/80">{title}</h3>
      {codes.map((referral, index) => (
        <ReferralCodeItem
          key={referral.id}
          referral={referral}
          index={index}
          isCopied={copiedIndex === index}
          onCopy={onCopy}
          onShare={onShare}
        />
      ))}
    </div>
  );
}

interface HowItWorksSectionProps {
  title: string;
  instructions: string[];
}

function HowItWorksSection({ title, instructions }: HowItWorksSectionProps) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-medium text-white/80">{title}</h3>
      <div className="p-3 rounded-lg bg-white/5 border border-white/10">
        <div className="text-xs sm:text-sm text-white/70 space-y-1">
          {instructions.map((instruction, index) => (
            <p key={index}>• {instruction}</p>
          ))}
        </div>
      </div>
    </div>
  );
}

// Main component
export function MyReferralModal({
  isOpen,
  onOpenChange,
}: MyReferralModalProps) {
  const { t } = useTranslation();
  const [copiedIndex, setCopiedIndex] = useState<number | null>(null);
  const [shareModalOpen, setShareModalOpen] = useState(false);
  const [selectedCode, setSelectedCode] = useState<string>("");
  const {
    data: referralCodesResponse,
    isLoading,
    error,
    refetch,
  } = useReferralCodes();
  const { setCodes, setLoading, setError } = useReferralStore();

  const referralCodes = referralCodesResponse?.data?.codes || [];

  // Sync data with store
  useEffect(() => {
    if (isOpen) {
      refetch();
    }
  }, [isOpen, refetch]);

  useEffect(() => {
    if (referralCodesResponse?.data?.codes) {
      setCodes(referralCodesResponse.data.codes);
    }
  }, [referralCodesResponse, setCodes]);

  useEffect(() => {
    setLoading(isLoading);
  }, [isLoading, setLoading]);

  useEffect(() => {
    if (error) {
      setError(error.message);
    }
  }, [error, setError]);

  // Handle copy action
  const handleCopy = useCallback(
    async (codeToCopy: string, index: number) => {
      try {
        await navigator.clipboard.writeText(codeToCopy);
        setCopiedIndex(index);

        setTimeout(() => setCopiedIndex(null), COPY_RESET_TIMEOUT);

        toast.success(t("waitlist.referral_code_copied"), TOAST_CONFIG.success);
      } catch (error) {
        logger.error("Failed to copy referral code:", error);
        toast.error(t("waitlist.failed_copy_code"), TOAST_CONFIG.error);
      }
    },
    [t]
  );

  // Handle share action
  const handleShare = useCallback((code: string) => {
    setSelectedCode(code);
    setShareModalOpen(true);
  }, []);

  const modalTitle = t("waitlist.your_referral_code");

  if (isLoading) {
    return (
      <ModalWrapper
        isOpen={isOpen}
        onOpenChange={onOpenChange}
        title={modalTitle}
      >
        <LoadingState />
      </ModalWrapper>
    );
  }

  if (error) {
    return (
      <ModalWrapper
        isOpen={isOpen}
        onOpenChange={onOpenChange}
        title={modalTitle}
      >
        <ErrorState
          message={t("waitlist.failed_load_codes")}
          errorDetails={error.message}
        />
      </ModalWrapper>
    );
  }

  return (
    <>
      <ModalWrapper
        isOpen={isOpen}
        onOpenChange={onOpenChange}
        title={modalTitle}
        description={t("waitlist.share_code_description")}
      >
        <div className="space-y-6 p-2">
          <ReferralCodesList
            codes={referralCodes}
            copiedIndex={copiedIndex}
            onCopy={handleCopy}
            onShare={handleShare}
            emptyMessage={t("waitlist.no_referral_codes")}
            title={t("waitlist.your_referral_codes")}
          />
          <HowItWorksSection
            title={t("waitlist.how_it_works")}
            instructions={[
              t("waitlist.share_code_instruction"),
              t("waitlist.skip_waitlist_instruction"),
            ]}
          />
        </div>
      </ModalWrapper>

      <ShareReferralModal
        isOpen={shareModalOpen}
        onOpenChange={setShareModalOpen}
        referralCode={selectedCode}
      />
    </>
  );
}

export default MyReferralModal;
