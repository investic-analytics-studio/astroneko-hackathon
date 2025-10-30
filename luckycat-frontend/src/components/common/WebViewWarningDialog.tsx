import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "../ui/dialog";
import { AlertTriangle } from "lucide-react";
import { ExternalBrowserButton } from "../ui/ExternalBrowserButton";
import {
  shouldShowExternalBrowserButton,
  isAndroid,
  isiOS,
} from "../../lib/webviewDetection";
import { useTranslation } from "react-i18next";

interface WebViewWarningDialogProps {
  readonly isOpen: boolean;
  readonly currentUrl: string;
}

export function WebViewWarningDialog({
  isOpen,
  currentUrl,
}: WebViewWarningDialogProps) {
  const { t } = useTranslation();

  const getButtonText = () => {
    if (isAndroid()) {
      return t("webview.open_in_browser");
    }
    if (isiOS()) {
      return t("webview.open_in_safari");
    }
    return t("webview.open_in_browser");
  };

  return (
    <Dialog open={isOpen} modal>
      <DialogContent
        showCloseButton={false}
        className="sm:max-w-md bg-black/95 backdrop-blur-xl border-white/20 rounded-[20px] text-white"
        onEscapeKeyDown={(e) => e.preventDefault()}
        onPointerDownOutside={(e) => e.preventDefault()}
        onInteractOutside={(e) => e.preventDefault()}
      >
        <DialogHeader className="text-center">
          <div className="mx-auto mb-2 flex h-10 w-10 items-center justify-center rounded-full bg-[#F7C36D]/20">
            <AlertTriangle className="h-5 w-5 text-[#F7C36D]" />
          </div>
          <DialogTitle className="text-lg font-bold text-[#F7C36D]">
            {t("webview.external_browser_required")}
          </DialogTitle>
          <DialogDescription className="text-[#A1A1AA] text-sm">
            {t("webview.google_login_security")}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {shouldShowExternalBrowserButton() && (
            <div className="flex justify-center">
              <ExternalBrowserButton
                targetUrl={currentUrl}
                variant="default"
                className="bg-[#F7C36D] text-black hover:bg-[#F7C36D]/90 font-medium"
              >
                {getButtonText()}
              </ExternalBrowserButton>
            </div>
          )}

          <p className="text-xs text-[#A1A1AA] text-center">
            {t("webview.dialog_security")}
          </p>
        </div>
      </DialogContent>
    </Dialog>
  );
}
