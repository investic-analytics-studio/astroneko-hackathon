import { useCallback, useState } from "react";

// Constants
const DEFAULT_RESET_DELAY_MS = 1500;
const FALLBACK_TEXTAREA_STYLE = {
  position: "fixed",
  left: "-999999px",
  top: "-999999px",
} as const;

// Type definitions
interface UseClipboardOptions {
  resetDelayMs?: number;
}

export interface UseClipboardReturn {
  copied: boolean;
  error: Error | null;
  copy: (value: string) => Promise<void>;
}

/**
 * Custom hook for clipboard operations
 * Provides a secure way to copy text with fallback support
 *
 * @param options - Configuration options
 * @returns Object with copied state, error, and copy function
 *
 * @example
 * ```tsx
 * const { copied, copy, error } = useClipboard();
 *
 * <button onClick={() => copy("Text to copy")}>
 *   {copied ? "Copied!" : "Copy"}
 * </button>
 * ```
 */
export const useClipboard = ({ resetDelayMs = DEFAULT_RESET_DELAY_MS }: UseClipboardOptions = {}): UseClipboardReturn => {
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const copy = useCallback(
    async (value: string) => {
      setError(null);

      try {
        // Use modern clipboard API if available
        if (navigator.clipboard && window.isSecureContext) {
          await navigator.clipboard.writeText(value);
        } else {
          // Fallback for older browsers or non-secure contexts
          const textArea = document.createElement("textarea");
          textArea.value = value;
          textArea.setAttribute("readonly", "");
          Object.assign(textArea.style, FALLBACK_TEXTAREA_STYLE);

          document.body.appendChild(textArea);
          textArea.focus();
          textArea.select();

          const successful = document.execCommand("copy");
          document.body.removeChild(textArea);

          if (!successful) {
            throw new Error("Fallback copy command failed");
          }
        }

        setCopied(true);
        setTimeout(() => setCopied(false), resetDelayMs);
      } catch (err) {
        const clipboardError = err instanceof Error ? err : new Error("Copy failed");
        setError(clipboardError);
        setCopied(false);
      }
    },
    [resetDelayMs]
  );

  return {
    copied,
    error,
    copy,
  };
};
