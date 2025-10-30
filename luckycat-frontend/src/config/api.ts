/**
 * API Configuration
 *
 * Centralized API configuration and constants
 */

/**
 * Base API URL from environment variables
 */
export const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || "";

/**
 * API version
 */
export const API_VERSION = "v1";

/**
 * API endpoints configuration
 */
export const API_ENDPOINTS = {
  auth: {
    base: `${API_BASE_URL}/${API_VERSION}/api/auth`,
    google: "/google",
    refresh: "/refresh",
    me: "/me",
    delete: "/delete",
    referral: {
      activate: "/referral/activate",
      codes: "/referral/codes",
    },
  },
  fortune: {
    base: `${API_BASE_URL}/${API_VERSION}/api/agent`,
    reply: "/reply",
    clearState: "/clear-state",
  },
  waitlist: {
    base: `${API_BASE_URL}/${API_VERSION}/api/waiting-list`,
    join: "/join",
    check: "/check",
  },
  userLimit: {
    base: `${API_BASE_URL}/${API_VERSION}/api/user-limit`,
    check: "/check",
  },
  chatHistory: {
    base: `${API_BASE_URL}/${API_VERSION}/api/history`,
    sessions: "/sessions",
    messages: "/sessions",
  },
} as const;

/**
 * API timeout configuration (in milliseconds)
 */
export const API_TIMEOUT = {
  default: 30000, // 30 seconds
  upload: 60000, // 60 seconds for file uploads
  download: 90000, // 90 seconds for downloads
} as const;

/**
 * API retry configuration
 */
export const API_RETRY = {
  attempts: 3,
  delay: 1000,
  backoff: 2, // Exponential backoff multiplier
} as const;

/**
 * HTTP status codes
 */
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  TOO_MANY_REQUESTS: 429,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
} as const;

/**
 * Helper to build full API URL
 */
export function buildApiUrl(endpoint: string): string {
  return `${API_BASE_URL}${endpoint}`;
}

/**
 * Helper to build query string from params
 */
export function buildQueryString(
  params: Record<string, string | number | boolean | undefined>
): string {
  const filtered = Object.entries(params)
    .filter(([_, value]) => value !== undefined && value !== null)
    .map(
      ([key, value]) =>
        `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`
    )
    .join("&");

  return filtered ? `?${filtered}` : "";
}
