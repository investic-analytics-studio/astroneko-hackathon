/**
 * API Error Handler Utilities
 *
 * Centralized error handling for API requests
 */

import { AxiosError } from "axios";
import { ApiError, type ApiErrorResponse } from "@/types/api";
import { HTTP_STATUS } from "@/config/api";

/**
 * Checks if error is an Axios error
 */
export function isAxiosError(error: unknown): error is AxiosError {
  return (
    typeof error === "object" &&
    error !== null &&
    "isAxiosError" in error &&
    error.isAxiosError === true
  );
}

/**
 * Extracts error message from API error response
 */
export function getErrorMessage(error: unknown): string {
  if (isAxiosError(error)) {
    const data = error.response?.data as ApiErrorResponse | undefined;

    // Check for structured error messages
    if (data?.status?.message && Array.isArray(data.status.message)) {
      return data.status.message.join(", ");
    }

    // Check for single error message
    if (data?.error) {
      return data.error;
    }

    // Check for validation errors
    if (data?.errors) {
      const messages = Object.values(data.errors).flat();
      return messages.join(", ");
    }

    // Fallback to axios error message
    return error.message || "An unexpected error occurred";
  }

  if (error instanceof Error) {
    return error.message;
  }

  return "An unexpected error occurred";
}

/**
 * Transforms axios error to ApiError
 */
export function handleApiError(error: unknown): never {
  if (isAxiosError(error)) {
    const data = error.response?.data as ApiErrorResponse | undefined;
    const statusCode = error.response?.status;
    const code = data?.status?.code;

    throw new ApiError(
      getErrorMessage(error),
      statusCode,
      code,
      data?.errors
    );
  }

  if (error instanceof Error) {
    throw new ApiError(error.message);
  }

  throw new ApiError("An unexpected error occurred");
}

/**
 * Check if error is a specific HTTP status
 */
export function isHttpStatus(error: unknown, status: number): boolean {
  return isAxiosError(error) && error.response?.status === status;
}

/**
 * Check if error is unauthorized (401)
 */
export function isUnauthorizedError(error: unknown): boolean {
  return isHttpStatus(error, HTTP_STATUS.UNAUTHORIZED);
}

/**
 * Check if error is forbidden (403)
 */
export function isForbiddenError(error: unknown): boolean {
  return isHttpStatus(error, HTTP_STATUS.FORBIDDEN);
}

/**
 * Check if error is not found (404)
 */
export function isNotFoundError(error: unknown): boolean {
  return isHttpStatus(error, HTTP_STATUS.NOT_FOUND);
}

/**
 * Check if error is rate limit (429)
 */
export function isRateLimitError(error: unknown): boolean {
  return isHttpStatus(error, HTTP_STATUS.TOO_MANY_REQUESTS);
}

/**
 * Check if error is server error (5xx)
 */
export function isServerError(error: unknown): boolean {
  if (!isAxiosError(error)) return false;
  const status = error.response?.status;
  return status !== undefined && status >= 500 && status < 600;
}

/**
 * Check if error is network error
 */
export function isNetworkError(error: unknown): boolean {
  return isAxiosError(error) && !error.response;
}

/**
 * Get user-friendly error message based on error type
 */
export function getUserFriendlyErrorMessage(error: unknown): string {
  if (isUnauthorizedError(error)) {
    return "Please login to continue";
  }

  if (isForbiddenError(error)) {
    return "You don't have permission to perform this action";
  }

  if (isNotFoundError(error)) {
    return "The requested resource was not found";
  }

  if (isRateLimitError(error)) {
    return "Too many requests. Please try again later";
  }

  if (isServerError(error)) {
    return "Server error. Please try again later";
  }

  if (isNetworkError(error)) {
    return "Network error. Please check your connection";
  }

  return getErrorMessage(error);
}