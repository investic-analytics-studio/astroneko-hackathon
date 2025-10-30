import * as React from "react";
import { AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";

export interface ErrorMessageProps extends React.HTMLAttributes<HTMLDivElement> {
  /**
   * Error message to display
   */
  message: string;
  /**
   * Optional error title
   */
  title?: string;
  /**
   * Show icon (default: true)
   */
  showIcon?: boolean;
  /**
   * Variant styling
   */
  variant?: "default" | "inline" | "card";
}

/**
 * ErrorMessage Component
 *
 * Consistent error message display with optional icon and variants.
 *
 * @example
 * ```tsx
 * // Basic usage
 * <ErrorMessage message="Something went wrong" />
 *
 * // With title
 * <ErrorMessage
 *   title="Error"
 *   message="Failed to load data"
 * />
 *
 * // Inline variant
 * <ErrorMessage
 *   variant="inline"
 *   message="Invalid email"
 * />
 *
 * // Card variant
 * <ErrorMessage
 *   variant="card"
 *   title="Upload Failed"
 *   message="File size exceeds limit"
 * />
 * ```
 */
export const ErrorMessage = React.forwardRef<HTMLDivElement, ErrorMessageProps>(
  (
    {
      message,
      title,
      showIcon = true,
      variant = "default",
      className,
      ...props
    },
    ref
  ) => {
    const variantClasses = {
      default: "text-red-500 text-sm",
      inline: "text-red-500 text-xs",
      card: "bg-red-500/10 border border-red-500/20 rounded-lg p-4 text-red-500",
    };

    return (
      <div
        ref={ref}
        role="alert"
        aria-live="polite"
        className={cn(variantClasses[variant], className)}
        {...props}
      >
        <div className="flex items-start gap-2">
          {showIcon && variant !== "inline" && (
            <AlertCircle className="w-4 h-4 mt-0.5 shrink-0" />
          )}
          <div className="flex-1">
            {title && variant === "card" && (
              <p className="font-semibold mb-1">{title}</p>
            )}
            <p>{message}</p>
          </div>
        </div>
      </div>
    );
  }
);

ErrorMessage.displayName = "ErrorMessage";