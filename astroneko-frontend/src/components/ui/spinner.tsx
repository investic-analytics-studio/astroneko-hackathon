import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

/**
 * Spinner variants using CVA for consistent styling
 */
const spinnerVariants = cva("animate-spin rounded-full border-solid", {
  variants: {
    size: {
      sm: "h-4 w-4 border-2",
      md: "h-6 w-6 border-2",
      lg: "h-8 w-8 border-3",
      xl: "h-12 w-12 border-4",
    },
    variant: {
      default: "border-white/20 border-t-white",
      primary: "border-[var(--brand-primary)]/20 border-t-[var(--brand-primary)]",
      accent: "border-[var(--brand-accent)]/20 border-t-[var(--brand-accent)]",
      dark: "border-black/20 border-t-black",
    },
  },
  defaultVariants: {
    size: "md",
    variant: "default",
  },
});

export interface SpinnerProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof spinnerVariants> {
  /**
   * Accessible label for screen readers
   */
  label?: string;
}

/**
 * Spinner Component
 *
 * A reusable loading spinner with multiple size and color variants.
 * Uses CSS variables from index.css for brand colors.
 *
 * @example
 * ```tsx
 * // Basic usage
 * <Spinner />
 *
 * // With variants
 * <Spinner size="lg" variant="primary" />
 *
 * // With custom label
 * <Spinner label="Loading data..." />
 * ```
 */
export const Spinner = React.forwardRef<HTMLDivElement, SpinnerProps>(
  ({ className, size, variant, label = "Loading...", ...props }, ref) => {
    return (
      <div
        ref={ref}
        role="status"
        aria-label={label}
        className={cn(spinnerVariants({ size, variant }), className)}
        {...props}
      >
        <span className="sr-only">{label}</span>
      </div>
    );
  }
);

Spinner.displayName = "Spinner";