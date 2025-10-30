import axios from "@/config/axios";
import { API_ENDPOINTS } from "@/config/api";
import type { ApiResponse } from "@/types/api";

export interface AuthUser {
  id: string;
  email: string;
  is_activated_referral: boolean;
  latest_login_at: string;
  firebase_uid: string;
  profile_image_url: string | null;
  display_name: string | null;
  created_at: string;
  updated_at: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface AuthResponse {
  status: {
    code: string;
    message: string[];
  };
  data: {
    user: AuthUser;
  } & AuthTokens;
}

export interface GoogleAuthRequest {
  id_token: string;
  refresh_token: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ReferralRequest {
  referral_code: string;
}

export interface ReferralCode {
  id: string;
  referral_code: string;
  is_activated: boolean;
}

export interface ReferralCodesResponse {
  status: {
    code: string;
    message: string[];
  };
  data: {
    codes: ReferralCode[];
  };
}

export interface ActivateReferralResponse {
  status: {
    code: string;
    message: string[];
  };
  data: {
    success: boolean;
    message: string;
  };
}

export async function authenticateWithGoogle(
  params: GoogleAuthRequest
): Promise<AuthResponse> {
  const response = await axios.post<AuthResponse>(
    `${API_ENDPOINTS.auth.base}${API_ENDPOINTS.auth.google}`,
    {
      id_token: params.id_token,
      refresh_token: params.refresh_token,
    }
  );

  return response.data;
}

export async function refreshToken(
  params: RefreshTokenRequest
): Promise<AuthResponse> {
  const response = await axios.post<AuthResponse>(
    `${API_ENDPOINTS.auth.base}${API_ENDPOINTS.auth.refresh}`,
    params
  );

  // Store the new tokens from refresh response
  // Using dynamic import to avoid circular dependency
  if (response.data?.data) {
    const { setAccessToken, setRefreshToken } = await import("@/lib/cookie");

    const {
      access_token,
      expires_in,
      refresh_token: newRefreshToken,
    } = response.data.data;

    if (access_token && expires_in) {
      setAccessToken(access_token, expires_in);
    }

    if (newRefreshToken) {
      setRefreshToken(newRefreshToken);
    }
  }

  return response.data;
}

export async function getMe(): Promise<ApiResponse<AuthUser>> {
  const response = await axios.get<ApiResponse<AuthUser>>(
    `${API_ENDPOINTS.auth.base}${API_ENDPOINTS.auth.me}`
  );

  return response.data;
}

export async function deleteUserFirebase(
  firebaseUid: string
): Promise<ApiResponse<unknown>> {
  const response = await axios.delete<ApiResponse<unknown>>(
    `${API_ENDPOINTS.auth.base}${API_ENDPOINTS.auth.delete}?firebase_uid=${firebaseUid}`
  );

  return response.data;
}

export async function activateReferral(
  params: ReferralRequest
): Promise<ActivateReferralResponse> {
  const response = await axios.post<ActivateReferralResponse>(
    `${API_ENDPOINTS.auth.base}${API_ENDPOINTS.auth.referral.activate}`,
    params
  );

  return response.data;
}

export async function getUserReferralCode(): Promise<ReferralCodesResponse> {
  const response = await axios.get<ReferralCodesResponse>(
    `${API_ENDPOINTS.auth.base}${API_ENDPOINTS.auth.referral.codes}`
  );

  return response.data;
}
