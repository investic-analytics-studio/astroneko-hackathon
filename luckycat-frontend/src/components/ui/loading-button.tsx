import * as React from "react";
import { Button, type ButtonProps } from "./button";
import { Spinner } from "./spinner";
import { cn } from "@/lib/utils";

export interface LoadingButtonProps extends ButtonProps {
  /**
   * Loading state - disables button and shows spinner
   */
  loading?: boolean;
  /**
   * Text to show when loading (optional)
   */
  loadingText?: string;
  /**
   * Spinner size (defaults to "sm")
   */
  spinnerSize?: "sm" | "md" | "lg" | "xl";
}

/**
 * LoadingButton Component
 *
 * Button component with integrated loading state and spinner.
 * Automatically disables interaction when loading.
 *
 * @example
 * ```tsx
 * // Basic usage
 * <LoadingButton loading={isSubmitting}>
 *   Submit
 * </LoadingButton>
 *
 * // With loading text
 * <LoadingButton loading={isSubmitting} loadingText="Submitting...">
 *   Submit Form
 * </LoadingButton>
 *
 * // With variants
 * <LoadingButton variant="outline" size="lg" loading={isLoading}>
 *   Click Me
 * </LoadingButton>
 * ```
 */
export const LoadingButton = React.forwardRef<
  HTMLButtonElement,
  LoadingButtonProps
>(
  (
    {
      children,
      loading = false,
      loadingText,
      spinnerSize = "sm",
      disabled,
      className,
      ...props
    },
    ref
  ) => {
    return (
      <Button
        ref={ref}
        disabled={disabled || loading}
        aria-busy={loading}
        aria-live={loading ? "polite" : undefined}
        className={cn(className)}
        {...props}
      >
        {loading ? (
          <div className="flex items-center gap-2">
            <Spinner size={spinnerSize} variant="default" aria-hidden="true" />
            {loadingText || children}
          </div>
        ) : (
          children
        )}
      </Button>
    );
  }
);

LoadingButton.displayName = "LoadingButton";