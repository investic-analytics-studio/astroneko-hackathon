import { useState } from "react";
import { AlertTriangle, Trash2, X } from "lucide-react";
import { useTranslation } from "react-i18next";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { logger } from "@/lib/logger";

interface DeleteSessionModalProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  sessionName: string;
  onConfirmDelete: () => Promise<void>;
  isDeleting?: boolean;
}

export const DeleteSessionModal = ({
  isOpen,
  onOpenChange,
  sessionName,
  onConfirmDelete,
  isDeleting = false,
}: DeleteSessionModalProps) => {
  const { t } = useTranslation();
  const [isProcessing, setIsProcessing] = useState(false);

  const handleDelete = async () => {
    setIsProcessing(true);
    try {
      await onConfirmDelete();
      onOpenChange(false);
    } catch (error) {
      logger.error("Failed to delete session:", error);
    } finally {
      setIsProcessing(false);
    }
  };

  const isLoading = isDeleting || isProcessing;

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[90vw] max-w-[400px] bg-black/95 backdrop-blur-xl border-red-500/30 rounded-[16px] sm:rounded-[20px] text-white p-3 sm:p-4 md:p-6 lg:p-8 max-h-[85vh] sm:max-h-[90vh] overflow-y-auto">
        <DialogHeader className="space-y-3 sm:space-y-4">
          <div className="flex items-center justify-center mb-2 sm:mb-4 animate-scale-in">
            <div className="relative">
              <div className="absolute inset-0 bg-red-500/20 rounded-full blur-xl animate-pulse-glow" />
              <div className="relative w-12 h-12 sm:w-16 sm:h-16 rounded-full bg-gradient-to-br from-red-500/30 to-red-600/30 border border-red-500/40 flex items-center justify-center">
                <AlertTriangle className="w-6 h-6 sm:w-8 sm:h-8 text-red-400" />
              </div>
            </div>
          </div>

          <DialogTitle className="text-lg sm:text-xl md:text-2xl font-bold text-center text-white leading-tight px-1">
            {t("chat.delete_conversation_title", "Delete Conversation")}
          </DialogTitle>

          <DialogDescription className="text-center text-gray-300 pt-1 sm:pt-2 text-sm sm:text-base leading-relaxed px-2 sm:px-4">
            {t(
              "chat.delete_conversation_description",
              "Are you sure you want to delete this conversation? This action cannot be undone."
            )}
          </DialogDescription>
        </DialogHeader>

        <div
          className="sm:mt-4 p-2 sm:p-3 md:p-4 rounded-lg bg-red-500/10 border border-red-500/20 opacity-0 translate-y-2 animate-fade-in"
          style={{ animationDelay: "0.1s" }}
        >
          <div className="flex items-start gap-2 sm:gap-3">
            <Trash2 className="w-4 h-4 sm:w-5 sm:h-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div className="flex-1 min-w-0">
              <p className="text-xs sm:text-sm font-medium text-red-300">
                {t("chat.conversation_to_delete", "Conversation to delete:")}
              </p>
              <p className="text-sm sm:text-base font-semibold text-white mt-1 break-words">
                {sessionName}
              </p>
            </div>
          </div>
        </div>

        <DialogFooter className="mt-3 sm:mt-4 md:mt-6 gap-3 flex-col sm:flex-row sm:justify-center">
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
            className="w-full sm:w-auto sm:min-w-1/2 h-9 sm:h-10 md:h-11 bg-white/10 hover:bg-white/20 border-white/20 text-white transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed order-2 sm:order-1 text-sm sm:text-base hover:text-white"
          >
            <X className="w-4 h-4 mr-2" />
            {t("common.cancel", "Cancel")}
          </Button>

          <Button
            type="button"
            onClick={handleDelete}
            disabled={isLoading}
            className="w-full sm:w-auto sm:min-w-1/2 h-9 sm:h-10 md:h-11 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white border-none shadow-lg shadow-red-500/20 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none order-1 sm:order-2 text-sm sm:text-base"
          >
            {isLoading ? (
              <>
                <div className="w-4 h-4 mr-2 border-2 border-white border-t-transparent rounded-full animate-spin" />
                {t("common.deleting", "Deleting...")}
              </>
            ) : (
              <>
                <Trash2 className="w-4 h-4 mr-2" />
                {t("common.delete", "Delete")}
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
