import ReactGA from "react-ga4";
import { v4 as uuidv4 } from "uuid";

enum HitTypes {
  PageView = "pageview",
}

type GtagFunction = (
  command: "event",
  eventName: string,
  params?: Record<string, unknown>
) => void;

declare global {
  interface Window {
    gtag?: GtagFunction;
  }
}

export const gaEnvKey = import.meta.env.VITE_GA_PROJECT_ID;

const GA_CLIENT_STORAGE_KEY = "ga-gtag-client-id";
const GA_DEVICE_ID_KEY = "ga-device-id";

let initialized = false;

/**
 * Get English category name for tracking events
 * @param categoryId - The category ID (e.g., 'general', 'crypto', 'lover', 'tarot')
 * @returns English category name for consistent tracking
 */
function getCategoryEnglishName(categoryId: string): string {
  const categoryNames: Record<string, string> = {
    general: "General",
    crypto: "Crypto",
    lover: "Lover",
    tarot: "Tarot",
  };

  return categoryNames[categoryId] || categoryId;
}

const generateClientIdGa = () => {
  let clientId = localStorage.getItem(GA_CLIENT_STORAGE_KEY);

  if (!clientId) {
    clientId = uuidv4();
    localStorage.setItem(GA_CLIENT_STORAGE_KEY, clientId);
  }

  return clientId;
};

const generateDeviceId = () => {
  let deviceId = localStorage.getItem(GA_DEVICE_ID_KEY);

  if (!deviceId) {
    deviceId = uuidv4();
    localStorage.setItem(GA_DEVICE_ID_KEY, deviceId);
  }

  return deviceId;
};

export interface CategorySelectedParams {
  category_id: string;
  category_name: string;
  category_type: string;
  category_subtitle?: string;
}

export interface Ga4CustomEventOptions {
  eventName: string;
  params?: {
    [key: string]: unknown;
  };
}

export function googleAnalyticsInit() {
  if (initialized) {
    console.debug("[Google Analytics] Already initialized");
    return;
  }

  const trackingId = gaEnvKey;
  if (!trackingId) {
    console.warn(
      "[Google Analytics] No tracking id found. Analytics disabled."
    );
    return;
  }

  try {
    ReactGA.initialize([
      {
        trackingId,
        gaOptions: {
          anonymizeIp: true,
          clientId: generateClientIdGa(),
          deviceId: generateDeviceId(),
        },
      },
    ]);

    initialized = true;

    // Set default user properties (similar to Amplitude's anonymous properties)
    setDefaultUserProperties();
  } catch (error) {
    console.error("[Google Analytics] Initialization failed:", error);
  }
}

/**
 * Set default user properties for all users (similar to Amplitude's anonymous properties)
 */
function setDefaultUserProperties(): void {
  try {
    const defaultProps = {
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

    setUserProperties(defaultProps);
  } catch (error) {
    console.warn("[Google Analytics] Failed to set default properties:", error);
  }
}

/**
 * Get browser name and version (same as Amplitude)
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
 * Set user identity after authentication (similar to Amplitude's setUser)
 * @param userId - User identifier
 */
export function setGaUser(userId: string | null): void {
  if (!initialized) return;

  try {
    if (userId) {
      ReactGA.set({ userId });
      console.debug("[Google Analytics] User ID set:", userId);
    } else {
      // Logout: clear user data
      ReactGA.set({ userId: undefined });
      console.debug("[Google Analytics] User cleared");
    }
  } catch (error) {
    console.error("[Google Analytics] setGaUser failed:", error);
  }
}

/**
 * Set custom user properties (similar to Amplitude's identifyUser)
 * @param props - User properties to set
 */
export function setUserProperties(props: Record<string, unknown>): void {
  if (!initialized) return;

  try {
    ReactGA.set(props);
  } catch (error) {
    console.error("[Google Analytics] setUserProperties failed:", error);
  }
}

/**
 * Check if Google Analytics is initialized and ready
 */
export function isGaInitialized(): boolean {
  return initialized;
}

/**
 * Track custom event with lowercase event names (identical to Amplitude)
 * @param event - Event name (lowercase, e.g., "launch app clicked")
 * @param props - Event properties (context-specific data)
 */
export function trackGaEvent(
  event: string,
  props?: Record<string, unknown>
): void {
  if (!initialized) {
    console.debug("[Google Analytics] Not initialized, skipping event:", event);
    return;
  }

  try {
    // Use the exact same lowercase event name as Amplitude
    const gaEventName = event.toLowerCase();

    // Add automatic context data similar to Amplitude
    const enrichedProps = {
      ...props,
      timestamp: new Date().toISOString(),
      referrer: document.referrer,
    };

    ReactGA.event(gaEventName, enrichedProps);
  } catch (error) {
    console.error("[Google Analytics] trackGaEvent failed:", error);
  }
}

/**
 * Track page view event (similar to Amplitude's trackPageView)
 * @param path - Current page path
 * @param search - Optional query string
 * @param pageName - Optional page name
 */
export function trackGaPageView(
  path: string,
  search?: string,
  pageName?: string
): void {
  if (!initialized) {
    console.debug("[Google Analytics] Not initialized, skipping page view");
    return;
  }

  try {
    // Also track page view as an event for consistency with Amplitude
    const eventName = pageName
      ? `page view - ${pageName.toLowerCase()}`
      : "page view";
    trackGaEvent(eventName, {
      path,
      search: search || "",
      page_name: pageName || path,
    });

    ReactGA.send({
      hitType: HitTypes.PageView,
      page: path + (search || ""),
      title: pageName || document.title,
    });
  } catch (error) {
    console.error("[Google Analytics] trackGaPageView failed:", error);
  }
}

// Legacy tracking functions - now use lowercase event names identical to Amplitude
export const trackCategorySelected = (id: string, subtitle?: string) => {
  const englishTitle = getCategoryEnglishName(id);
  trackGaEvent(`category selected - ${englishTitle.toLowerCase()}`, {
    category_id: id,
    category_name: englishTitle,
    category_type: id.toLowerCase(),
    category_subtitle: subtitle,
  });
};

export const trackWaitlistClick = (placement: string, buttonText: string) => {
  trackGaEvent("waitlist clicked", {
    placement,
    button_text: buttonText,
  });
};

export const trackLaunchAppClick = (placement: string, buttonText: string) => {
  trackGaEvent("launch app clicked", {
    placement,
    button_text: buttonText,
  });
};

export const trackLoginClick = (placement: string, buttonText: string) => {
  trackGaEvent("login clicked", {
    placement,
    button_text: buttonText,
  });
};

// Add tracking function for Referral modal clicked (matching Amplitude)
export const trackReferralModalClick = (
  placement: string,
  buttonText: string
) => {
  trackGaEvent("referral modal clicked", {
    placement,
    button_text: buttonText,
  });
};

// Add tracking function for Your Referral Code button clicked
export const trackYourReferralCodeClick = (
  placement: string,
  buttonText: string
) => {
  trackGaEvent("your referral code clicked", {
    placement,
    button_text: buttonText,
  });
};

// Type-safe event tracking helpers (similar to Amplitude's event helpers)
export const gaAuthEvents = {
  loginAttempt: (method: "email" | "google") =>
    trackGaEvent("login attempted", { method }),

  loginSuccess: (method: "email" | "google") =>
    trackGaEvent("login success", { method }),

  loginFailure: (method: "email" | "google", error: string) =>
    trackGaEvent("login failed", { method, error }),

  logout: () => trackGaEvent("logout"),
} as const;

export const gaFeatureEvents = {
  used: (
    feature: string,
    result: "success" | "error",
    metadata?: Record<string, unknown>
  ) => trackGaEvent("feature used", { feature, result, ...metadata }),
} as const;

export const gaErrorEvents = {
  occurred: (
    code: string,
    message: string,
    area: string,
    metadata?: Record<string, unknown>
  ) => trackGaEvent("error occurred", { code, message, area, ...metadata }),
} as const;

/**
 * Capture UTM parameters from URL and save to user properties.
 * Similar to Amplitude's captureUTMParams.
 * @param searchParams - URLSearchParams or search string
 */
export function captureGaUTMParams(
  searchParams: URLSearchParams | string
): void {
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
    setUserProperties(utmParams);
    console.debug("[Google Analytics] UTM parameters captured:", utmParams);
  }
}

export default {
  googleAnalyticsInit,
  setOption: (key: string, value: unknown) => ReactGA.set({ [key]: value }),
  setUserId: setGaUser,
  setUserProperties,
  sendData: (type: HitTypes, data: object) =>
    ReactGA.send({ hitType: type, ...data }),
  trackPageView: trackGaPageView,
  trackEvent: trackGaEvent,
  trackCategorySelected,
  trackWaitlistClick,
  trackLaunchAppClick,
  trackLoginClick,
  trackReferralModalClick,
  trackYourReferralCodeClick,
  authEvents: gaAuthEvents,
  featureEvents: gaFeatureEvents,
  errorEvents: gaErrorEvents,
  captureUTMParams: captureGaUTMParams,
  isInitialized: isGaInitialized,
};
