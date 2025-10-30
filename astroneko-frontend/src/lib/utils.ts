import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { toast } from "sonner";
import { logger } from "./logger";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getImagePath(path: string): string {
  const cleanPath = path.startsWith("/") ? path.slice(1) : path;
  return `/${cleanPath}`;
}

export function preloadImage(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onload = () => resolve();
    img.onerror = () => reject(new Error(`Failed to load image: ${src}`));
    img.src = getImagePath(src);
  });
}

export function preloadImages(imagePaths: string[]): Promise<void[]> {
  return Promise.all(imagePaths.map(preloadImage));
}

export function isMobileDevice(): boolean {
  if (typeof window === "undefined") return false;

  return (
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
      navigator.userAgent
    ) || window.innerWidth <= 768
  );
}

export async function copyToClipboard(
  text: string,
  successMessage = "Copied to clipboard!",
  errorMessage = "Failed to copy to clipboard"
): Promise<void> {
  try {
    await navigator.clipboard.writeText(text);
    toast.success(successMessage, {
      className: "custom-success-toast",
      descriptionClassName: "text-[#A1A1AA]",
    });
  } catch (error) {
    logger.error("Copy to clipboard failed:", error);
    toast.error(errorMessage, {
      className: "custom-error-toast",
      descriptionClassName: "text-[#A1A1AA]",
    });
    throw error;
  }
}

export function formatDate(
  date: Date | string,
  locale = "en-US",
  options: Intl.DateTimeFormatOptions = {
    year: "numeric",
    month: "long",
    day: "numeric",
  }
): string {
  const dateObj = typeof date === "string" ? new Date(date) : date;
  return new Intl.DateTimeFormat(locale, options).format(dateObj);
}

export function truncateText(
  text: string,
  maxLength: number,
  suffix = "..."
): string {
  if (text.length <= maxLength) return text;
  return `${text.slice(0, maxLength - suffix.length)}${suffix}`;
}

export function debounce<T extends (...args: unknown[]) => unknown>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout | null = null;

  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      timeout = null;
      func(...args);
    };

    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

export function formatNumber(
  value: number,
  locale = "en-US",
  options?: Intl.NumberFormatOptions
): string {
  return new Intl.NumberFormat(locale, options).format(value);
}

export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function safeJsonParse<T>(json: string, fallback: T): T {
  try {
    return JSON.parse(json) as T;
  } catch {
    return fallback;
  }
}
