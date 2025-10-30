/**
 * Shared Type Definitions
 *
 * Common types used across the application
 */

// Re-export specific types
export type { ApiResponse, ApiError, PaginationParams, PaginatedResponse } from './api';
export type { Card } from './pickCard';

/**
 * Common utility types
 */

/**
 * Make all properties optional recursively
 */
export type DeepPartial<T> = {
  [P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

/**
 * Make specific keys required in a type
 */
export type RequiredKeys<T, K extends keyof T> = T & Required<Pick<T, K>>;

/**
 * Make specific keys optional in a type
 */
export type OptionalKeys<T, K extends keyof T> = Omit<T, K> & Partial<Pick<T, K>>;

/**
 * Extract keys of type T that have value type V
 */
export type KeysOfType<T, V> = {
  [K in keyof T]: T[K] extends V ? K : never;
}[keyof T];

/**
 * Note: PaginationParams and PaginatedResponse are now imported from './api'
 * to avoid duplication and maintain consistency across the application
 */

/**
 * Loading state
 */
export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

/**
 * Async data state
 */
export interface AsyncData<T, E = Error> {
  data: T | null;
  error: E | null;
  isLoading: boolean;
  isError: boolean;
  isSuccess: boolean;
}

/**
 * Common callback types
 */
export type VoidCallback = () => void;
export type AsyncVoidCallback = () => Promise<void>;
export type ErrorCallback = (error: Error) => void;