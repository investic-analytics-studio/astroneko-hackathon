import * as React from "react";
import { Check, Copy } from "lucide-react";
import { Button, type ButtonProps } from "./button";
import { copyToClipboard } from "@/lib/utils";
import { cn } from "@/lib/utils";
import { logger } from "@/lib/logger";

export interface CopyButtonProps extends Omit<ButtonProps, "onClick"> {
  /**
   * Text to copy to clipboard
   */
  textToCopy: string;
  /**
   * Success message for toast (optional)
   */
  successMessage?: string;
  /**
   * Error message for toast (optional)
   */
  errorMessage?: string;
  /**
   * Duration to show success state in ms (default: 2000)
   */
  successDuration?: number;
  /**
   * Show text label (default: false)
   */
  showLabel?: boolean;
  /**
   * Custom label text
   */
  label?: string;
}

/**
 * CopyButton Component
 *
 * Button that copies text to clipboard with visual feedback.
 * Shows checkmark icon temporarily after successful copy.
 *
 * @example
 * ```tsx
 * // Basic usage
 * <CopyButton textToCopy="Hello World" />
 *
 * // With custom messages
 * <CopyButton
 *   textToCopy={referralCode}
 *   successMessage="Referral code copied!"
 *   errorMessage="Failed to copy code"
 * />
 *
 * // With label
 * <CopyButton
 *   textToCopy="https://example.com"
 *   showLabel
 *   label="Copy Link"
 *   variant="outline"
 * />
 *
 * // Icon only button
 * <CopyButton
 *   textToCopy="code"
 *   variant="ghost"
 *   size="sm"
 * />
 * ```
 */
export const CopyButton = React.forwardRef<HTMLButtonElement, CopyButtonProps>(
  (
    {
      textToCopy,
      successMessage,
      errorMessage,
      successDuration = 2000,
      showLabel = false,
      label = "Copy",
      className,
      children,
      ...props
    },
    ref
  ) => {
    const [copied, setCopied] = React.useState(false);

    const handleCopy = async () => {
      try {
        await copyToClipboard(textToCopy, successMessage, errorMessage);
        setCopied(true);

        // Reset copied state after duration
        setTimeout(() => {
          setCopied(false);
        }, successDuration);
      } catch (error) {
        // Error already handled in copyToClipboard utility
        logger.error("Copy failed:", error);
      }
    };

    return (
      <Button
        ref={ref}
        onClick={handleCopy}
        className={cn(className)}
        {...props}
      >
        {copied ? (
          <Check className="w-4 h-4" />
        ) : (
          <Copy className="w-4 h-4" />
        )}
        {showLabel && <span className="ml-2">{copied ? "Copied!" : label}</span>}
        {children}
      </Button>
    );
  }
);

CopyButton.displayName = "CopyButton";