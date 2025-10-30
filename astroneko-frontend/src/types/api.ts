/**
 * Shared API Types
 *
 * Common types used across API calls for consistent response handling
 */

/**
 * Standard API response wrapper
 */
export interface ApiResponse<T = unknown> {
  status: {
    code: string;
    message: string[];
  };
  data: T;
}

/**
 * API error response structure
 */
export interface ApiErrorResponse {
  status: {
    code: string;
    message: string[];
  };
  error?: string;
  errors?: Record<string, string[]>;
}

/**
 * Custom API Error class for better error handling
 */
export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public code?: string,
    public errors?: Record<string, string[]>
  ) {
    super(message);
    this.name = "ApiError";
    Object.setPrototypeOf(this, ApiError.prototype);
  }
}

/**
 * Pagination params for list endpoints
 */
export interface PaginationParams {
  page?: number;
  limit?: number;
  sort?: string;
  order?: "asc" | "desc";
}

/**
 * Paginated response wrapper
 */
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

/**
 * Common query parameters
 */
export interface QueryParams {
  [key: string]: string | number | boolean | undefined;
}