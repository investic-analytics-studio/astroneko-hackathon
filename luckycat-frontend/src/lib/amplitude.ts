/**
 * Amplitude Analytics Module
 *
 * Production-ready analytics wrapper for tracking user behavior and events.
 * Supports user identification, event tracking, and page view monitoring.
 *
 * @module lib/amplitude
 */

import {
  init as ampInit,
  track as ampTrack,
  identify as ampIdentify,
  Identify,
  setUserId as ampSetUserId,
  setDeviceId as ampSetDeviceId,
  reset as ampReset,
} from "@amplitude/analytics-browser";
import { logger } from "./logger";

// ---- Environment Configuration ----
const API_KEY = import.meta.env.VITE_AMPLITUDE_API_KEY || "";
const IS_DEV = import.meta.env.DEV;
const IS_PROD = import.meta.env.PROD;

let initialized = false;

/**
 * Initialize Amplitude SDK with production-ready configuration.
 * Should be called once during app bootstrap (client-side only).
 *
 * Features:
 * - Manual page view tracking (aligned with router)
 * - Automatic session tracking
 * - Debug logging in development
 * - Graceful fallback if API key is missing
 */
export function amplitudeInit(): void {
  if (initialized) {
    // logger.debug("[Amplitude] Already initialized");
    return;
  }

  if (!API_KEY) {
    logger.warn("[Amplitude] API key not found. Analytics disabled.");
    return;
  }

  try {
    ampInit(API_KEY, undefined, {
      // Manual tracking for better control
      defaultTracking: {
        pageViews: false, // We track manually via router
        sessions: true, // Auto session management
        formInteractions: false,
        fileDownloads: false,
      },
      autocapture: false,

      // Logging configuration
      logLevel: (IS_DEV ? 2 : 1) as 0 | 1 | 2 | 3 | 4, // Info (2) in dev, Warn (1) in prod

      // Optional: EU data residency
      // serverUrl: 'https://api.eu.amplitude.com/2/httpapi',

      // Performance optimization
      flushIntervalMillis: 1000,
      flushQueueSize: 30,
    });

    initialized = true;

    // Set anonymous user properties (for users who haven't logged in)
    setAnonymousUserProperties();
  } catch (error) {
    logger.error("[Amplitude] Initialization failed:", error);
  }
}

/**
 * Set properties for anonymous users (before login).
 * Helps understand who's using the app even without authentication.
 */
function setAnonymousUserProperties(): void {
  try {
    const anonymousProps: Record<string, unknown> = {
      // Browser & Device info
      browser: getBrowserInfo(),
      screen_width: window.screen.width,
      screen_height: window.screen.height,
      viewport_width: window.innerWidth,
      viewport_height: window.innerHeight,

      // Language & Location
      language: navigator.language,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,

      // Device capabilities
      is_mobile: /Mobi|Android/i.test(navigator.userAgent),
      is_tablet:
        /iPad|Android/i.test(navigator.userAgent) &&
        !/Mobile/i.test(navigator.userAgent),
      touch_support: "ontouchstart" in window || navigator.maxTouchPoints > 0,

      // Initial visit info
      initial_referrer: document.referrer || "direct",
      initial_landing_page: window.location.pathname,
    };

    identifyUser(anonymousProps);
  } catch (error) {
    logger.warn("[Amplitude] Failed to set anonymous properties:", error);
  }
}

/**
 * Get browser name and version
 */
function getBrowserInfo(): string {
  const ua = navigator.userAgent;
  if (ua.includes("Chrome")) return "Chrome";
  if (ua.includes("Safari")) return "Safari";
  if (ua.includes("Firefox")) return "Firefox";
  if (ua.includes("Edge")) return "Edge";
  return "Unknown";
}

/**
 * Set user identity after authentication.
 * Call this after successful login/signup.
 *
 * @param id - User identifier (preferably hashed for privacy)
 *
 * @example
 * ```ts
 * // After login
 * setUser(user.uid);
 * ```
 */
export function setUser(id: string | null): void {
  if (!initialized) return;

  try {
    if (id) {
      ampSetUserId(id);
      // logger.debug("[Amplitude] User ID set:", id);
    } else {
      // Logout: clear all user data
      ampReset();
      logger.debug("[Amplitude] User session reset");
    }
  } catch (error) {
    logger.error("[Amplitude] setUser failed:", error);
  }
}

/**
 * Set device identifier for cross-session tracking.
 *
 * @param deviceId - Unique device identifier
 */
export function setDevice(deviceId: string): void {
  if (!initialized) return;

  try {
    ampSetDeviceId(deviceId);
    logger.debug("[Amplitude] Device ID set:", deviceId);
  } catch (error) {
    logger.error("[Amplitude] setDevice failed:", error);
  }
}

/**
 * Identify user with custom properties.
 * Use for storing stable user attributes (plan, locale, signup_source, etc).
 *
 * @param props - User properties to set
 *
 * @example
 * ```ts
 * identifyUser({
 *   plan: 'pro',
 *   locale: 'th',
 *   signup_source: 'landing',
 * });
 * ```
 */
export function identifyUser(props: Record<string, unknown>): void {
  if (!initialized) return;

  try {
    const id = new Identify();
    Object.entries(props).forEach(([key, value]) => {
      id.set(key, value as string | number | boolean | string[] | number[]);
    });
    ampIdentify(id);
  } catch (error) {
    logger.error("[Amplitude] identifyUser failed:", error);
  }
}

/**
 * Track custom event with optional properties.
 * Use lowercase naming convention for consistency with GA.
 *
 * @param event - Event name (e.g., "cta clicked", "feature used")
 * @param props - Event properties (context-specific data)
 * @returns Promise that resolves when event is queued
 *
 * @example
 * ```ts
 * track('cta clicked', {
 *   cta: 'Start Now',
 *   placement: 'hero'
 * });
 *
 * track('feature used', {
 *   feature: 'export_csv',
 *   result: 'success'
 * });
 * ```
 */
export function track(event: string, props?: Record<string, unknown>): void {
  if (!initialized) {
    logger.debug("[Amplitude] Not initialized, skipping event:", event);
    return;
  }

  try {
    // Convert to lowercase for consistency with GA
    const eventName = event.toLowerCase();
    ampTrack(eventName, props);
  } catch (error) {
    logger.error("[Amplitude] track failed:", error);
  }
}

/**
 * Track page view event.
 * Typically called automatically by the router integration.
 *
 * @param path - Current page path
 * @param search - Optional query string
 *
 * @example
 * ```ts
 * trackPageView('/category/tarot', '?utm_source=google');
 * ```
 */
export function trackPageView(
  path: string,
  search?: string,
  pageName?: string
): void {
  const eventName = pageName
    ? `page view - ${pageName.toLowerCase()}`
    : "page view";
  track(eventName, {
    path,
    search: search || "",
    page_name: pageName || path,
    // Automatically add referrer and timestamp
    referrer: document.referrer,
    timestamp: new Date().toISOString(),
  });
}

/**
 * Capture UTM parameters from URL and save to user properties.
 * Call this on app init or first page load after consent.
 *
 * @param searchParams - URLSearchParams or search string
 *
 * @example
 * ```ts
 * // On first load
 * captureUTMParams(window.location.search);
 * ```
 */
export function captureUTMParams(searchParams: URLSearchParams | string): void {
  const params =
    typeof searchParams === "string"
      ? new URLSearchParams(searchParams)
      : searchParams;

  const utmParams = [
    "utm_source",
    "utm_medium",
    "utm_campaign",
    "utm_term",
    "utm_content",
  ].reduce((acc, key) => {
    const value = params.get(key);
    if (value) acc[key] = value;
    return acc;
  }, {} as Record<string, string>);

  if (Object.keys(utmParams).length > 0) {
    identifyUser(utmParams);
    logger.debug("[Amplitude] UTM parameters captured:", utmParams);
  }
}

// ---- Type-safe event tracking helpers ----

/**
 * Track authentication events
 */
export const authEvents = {
  loginAttempt: (method: "email" | "google") =>
    track("login attempted", { method }),

  loginSuccess: (method: "email" | "google") =>
    track("login success", { method }),

  loginFailure: (method: "email" | "google", error: string) =>
    track("login failed", { method, error }),

  logout: () => track("logout"),
} as const;

/**
 * Track feature usage
 */
export const featureEvents = {
  used: (
    feature: string,
    result: "success" | "error",
    metadata?: Record<string, unknown>
  ) => track("feature used", { feature, result, ...metadata }),
} as const;

/**
 * Track errors and exceptions
 */
export const errorEvents = {
  occurred: (
    code: string,
    message: string,
    area: string,
    metadata?: Record<string, unknown>
  ) => track("error occurred", { code, message, area, ...metadata }),
} as const;

/**
 * Check if Amplitude is initialized and ready
 */
export function isInitialized(): boolean {
  return initialized;
}

/**
 * Get initialization status for debugging
 */
export function getStatus(): {
  initialized: boolean;
  hasApiKey: boolean;
  env: string;
} {
  return {
    initialized,
    hasApiKey: !!API_KEY,
    env: IS_PROD ? "production" : "development",
  };
}
