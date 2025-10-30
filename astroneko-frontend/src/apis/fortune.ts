import { getSessionId, removeSessionId, setSessionId } from "@/lib/cookie";
import { logger } from "@/lib/logger";
import { API_ENDPOINTS } from "@/config/api";
import type { ApiResponse } from "@/types/api";
import axios from "../config/axios";

export interface FortuneResponse {
  status: string;
  message: string;
  card?: string;
  meaning?: string;
  session_id: string;
}

export interface FortuneMeta {
  is_guest: boolean;
  message: string;
}

export interface DailyLimitError {
  error: string;
  limit: number;
  message: string;
  reset_hours: number;
  reset_in: string;
  reset_mins: number;
  used: number;
}

export type FortuneErrorType =
  | "bad_request"
  | "free_trial_limit"
  | "daily_limit"
  | "daily_limit_logged_in"
  | "rate_limit"
  | "server_error"
  | "service_unavailable"
  | "network"
  | "unknown";

export interface LimitInfo {
  limit: number;
  used: number;
  resetHours: number;
  resetMins: number;
  resetIn: string;
}

export class FortuneError extends Error {
  constructor(
    public type: FortuneErrorType,
    message: string,
    public originalError: unknown,
    public limitInfo?: LimitInfo
  ) {
    super(message);
    this.name = "FortuneError";
  }
}

function createFortuneError(
  type: FortuneErrorType,
  message: string,
  originalError: unknown,
  limitInfo?: LimitInfo
): FortuneError {
  return new FortuneError(type, message, originalError, limitInfo);
}

function extractErrorData(error: unknown) {
  if (!error || typeof error !== "object" || !("response" in error)) {
    return null;
  }

  const axiosError = error as {
    response?: {
      status?: number;
      data?: Partial<DailyLimitError> & {
        error?: string;
        status?: {
          code?: string;
          message?: string[];
        };
      };
    };
  };

  return {
    data: axiosError.response?.data,
    statusCode: axiosError.response?.status,
  };
}

function handleLimitError(
  errorData: Partial<DailyLimitError>
): FortuneError {
  // Handle Free trial limit exceeded
  if (errorData.error === "Free trial limit exceeded") {
    const message = errorData.message || "You've used all 3 free requests. Please sign in for unlimited access.";
    const limitInfo =
      errorData.limit !== undefined && errorData.used !== undefined
        ? {
            limit: errorData.limit,
            used: errorData.used,
            resetHours: 0,
            resetMins: 0,
            resetIn: "",
          }
        : undefined;
    return createFortuneError("free_trial_limit", message, errorData, limitInfo);
  }

  // Handle Daily limit exceeded
  const isLoggedIn = errorData.message?.includes("activate a referral code");
  const type = isLoggedIn ? "daily_limit_logged_in" : "daily_limit";
  const message = isLoggedIn
    ? "You've used all your free daily requests. Please activate a referral code for unlimited access or try again tomorrow."
    : "You've used all your free daily requests. Please sign in for unlimited access or try again tomorrow.";

  const limitInfo =
    errorData.limit !== undefined && errorData.used !== undefined
      ? {
          limit: errorData.limit,
          used: errorData.used,
          resetHours: errorData.reset_hours || 0,
          resetMins: errorData.reset_mins || 0,
          resetIn: errorData.reset_in || "",
        }
      : undefined;

  return createFortuneError(type, message, errorData, limitInfo);
}

function handleStatusMessageError(
  statusMessage: string | string[] | undefined,
  defaultMessage: string
): string {
  if (!statusMessage) return defaultMessage;
  return Array.isArray(statusMessage)
    ? statusMessage.join(", ")
    : statusMessage;
}

function processApiError(error: unknown): FortuneError {
  const errorData = extractErrorData(error);

  if (!errorData) {
    return createFortuneError(
      "network",
      "Network error. Please check your connection and try again.",
      error
    );
  }

  const { data, statusCode } = errorData;

  switch (statusCode) {
    case 400:
      return createFortuneError(
        "bad_request",
        handleStatusMessageError(
          data?.status?.message,
          "Invalid request. Please check your input and try again."
        ),
        error
      );

    case 429:
      if (data?.error === "Daily limit exceeded" || data?.error === "Free trial limit exceeded") {
        return handleLimitError(data);
      } else if (data?.error === "Rate limit exceeded") {
        return createFortuneError(
          "rate_limit",
          "Our fortune teller is very popular right now! Please wait a few minutes and try again. 🐱",
          error
        );
      }
      break;

    case 500:
      return createFortuneError(
        "server_error",
        handleStatusMessageError(
          data?.status?.message,
          "Our fortune teller is experiencing issues. Please try again later."
        ),
        error
      );

    case 502:
      return createFortuneError(
        "service_unavailable",
        handleStatusMessageError(
          data?.status?.message,
          "Service temporarily unavailable. Please try again in a few moments."
        ),
        error
      );
  }

  return createFortuneError("unknown", "An unexpected error occurred", error);
}

export async function getFortuneReply(
  message: string,
  category: string = "general",
  chatHistorySessionId?: string | null
): Promise<FortuneResponse> {
  const sessionId = chatHistorySessionId || getSessionId(category);

  try {
    const response = await axios.post<ApiResponse<FortuneResponse>>(
      `${API_ENDPOINTS.fortune.base}${API_ENDPOINTS.fortune.reply}`,
      {
        text: message,
        session_id: sessionId ?? null,
      }
    );

    const fortuneData = response.data.data;
    if (fortuneData.session_id && !chatHistorySessionId) {
      setSessionId(fortuneData.session_id, category);
    }

    return fortuneData;
  } catch (error: unknown) {
    throw processApiError(error);
  }
}

export const FORTUNE_CATEGORIES = ["general", "crypto", "lover"] as const;

export type FortuneCategory = (typeof FORTUNE_CATEGORIES)[number];

function isValidCategory(category: string): category is FortuneCategory {
  return FORTUNE_CATEGORIES.includes(category as FortuneCategory);
}

export async function clearFortuneState(
  category: FortuneCategory = "general"
): Promise<void> {
  if (!isValidCategory(category)) {
    throw new Error(`Invalid fortune category: ${category}`);
  }

  const sessionId = getSessionId(category);
  if (!sessionId) {
    removeSessionId(category);
    return;
  }

  try {
    await axios.delete<ApiResponse<{ status: string }>>(
      `${API_ENDPOINTS.fortune.base}${API_ENDPOINTS.fortune.clearState}`,
      {
        data: { session_id: sessionId },
      }
    );
  } catch (error) {
    logger.warn(`Failed to clear backend session for ${category}:`, error);
  } finally {
    removeSessionId(category);
  }
}

export async function clearAllFortuneStates(): Promise<void> {
  const clearPromises = FORTUNE_CATEGORIES.map(async (category) => {
    try {
      await clearFortuneState(category);
    } catch (error) {
      logger.error(
        `Failed to clear fortune state for category ${category}:`,
        error
      );
    }
  });

  await Promise.allSettled(clearPromises);
}

export function isFortuneError(error: unknown): error is FortuneError {
  return error instanceof FortuneError;
}

export function getFortuneErrorMessage(error: unknown): string {
  if (isFortuneError(error)) {
    return error.message;
  }
  return "An unexpected error occurred. Please try again.";
}

export function isRateLimitError(error: unknown): boolean {
  return (
    isFortuneError(error) &&
    (error.type === "rate_limit" ||
      error.type === "free_trial_limit" ||
      error.type === "daily_limit" ||
      error.type === "daily_limit_logged_in")
  );
}

export function getDailyLimitResetTime(error: unknown): string | null {
  if (isFortuneError(error) && error.limitInfo) {
    return error.limitInfo.resetIn;
  }
  return null;
}
