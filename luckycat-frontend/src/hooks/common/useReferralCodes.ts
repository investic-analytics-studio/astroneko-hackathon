import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getUserReferralCode, ReferralCodesResponse } from "@/apis/auth";
import { useAuth } from "./useAuth";

const REFERRAL_CODES_QUERY_KEY = ["referral-codes"];

export const useReferralCodes = () => {
  const { isAuthenticated } = useAuth();

  return useQuery<ReferralCodesResponse, Error>({
    queryKey: REFERRAL_CODES_QUERY_KEY,
    queryFn: getUserReferralCode,
    enabled: isAuthenticated, // Only fetch when user is authenticated
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    retry: 1,
    refetchOnWindowFocus: false,
  });
};

export const useInvalidateReferralCodes = () => {
  const queryClient = useQueryClient();

  return () => {
    queryClient.invalidateQueries({
      queryKey: REFERRAL_CODES_QUERY_KEY,
    });
  };
};
