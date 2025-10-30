import { useState, useEffect, useCallback } from "react";
import { isUserLimitReached } from "@/apis/user-limit";
import { useAuth } from "./useAuth";
import { logger } from "@/lib/logger";

// Type definitions
export interface UseUserLimitReturn {
  isOverLimit: boolean | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

/**
 * Custom hook for checking user usage limits
 * Monitors if user has exceeded their usage quota
 *
 * @returns Object with limit state, loading state, error, and refetch function
 *
 * @example
 * ```tsx
 * const { isOverLimit, isLoading, refetch } = useUserLimit();
 *
 * if (isOverLimit) {
 *   return <LimitReachedMessage />;
 * }
 * ```
 */
export const useUserLimit = (): UseUserLimitReturn => {
  const [isOverLimit, setIsOverLimit] = useState<boolean | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { authUser } = useAuth();

  const checkUserLimit = useCallback(async () => {
    if (!authUser) {
      setIsOverLimit(null);
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError(null);
      const response = await isUserLimitReached();
      setIsOverLimit(response.data.is_limit_reached);
    } catch (err) {
      logger.error("Failed to check user limit:", err);
      setError("Failed to check user limit");
      setIsOverLimit(null);
    } finally {
      setIsLoading(false);
    }
  }, [authUser]);

  useEffect(() => {
    checkUserLimit();
  }, [authUser, checkUserLimit]);

  return {
    isOverLimit,
    isLoading,
    error,
    refetch: checkUserLimit,
  };
};
