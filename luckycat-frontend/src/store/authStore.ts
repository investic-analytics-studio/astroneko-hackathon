import { create } from "zustand";
import {
  signInWithPopup,
  signInWithRedirect,
  getRedirectResult,
  signOut,
  User as FirebaseUser,
} from "firebase/auth";
import { auth, googleProvider } from "../config/firebase";
import { logger } from "../lib/logger";
import {
  authenticateWithGoogle,
  getMe,
  type AuthUser,
  type AuthResponse,
} from "../apis/auth";
import {
  setAccessToken,
  setRefreshToken,
  clearAuthTokens,
  getAccessToken,
  getRefreshToken,
} from "../lib/cookie";
import { setUser, identifyUser, authEvents } from "../lib/amplitude";

interface AuthState {
  googleUser: FirebaseUser | null;
  authUser: AuthUser | null;
  isLoading: boolean;
  error: string | null;
  isInitializing: boolean; // Add flag to prevent duplicate initialization
}

interface AuthActions {
  initializeAuth: () => Promise<void>;
  loginWithGoogle: () => Promise<AuthResponse | null>;
  loginWithGoogleRedirect: () => Promise<void>;
  handleRedirectResult: () => Promise<AuthResponse | null>;
  logout: () => Promise<void>;
  refreshUserInfo: () => Promise<AuthUser | null>;
  checkAuthentication: () => Promise<boolean>;
  setGoogleUser: (user: FirebaseUser | null) => void;
  setAuthUser: (user: AuthUser | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  resetLoginState: () => void;
  clearAuthState: () => void;
}

type AuthStore = AuthState & AuthActions;

const initialState: AuthState = {
  googleUser: null,
  authUser: null,
  isLoading: true,
  error: null,
  isInitializing: false,
};

const handleAuthentication = async (
  set: (partial: Partial<AuthState>) => void
): Promise<boolean> => {
  const accessToken = getAccessToken();
  const refreshTokenValue = getRefreshToken();

  if (accessToken) {
    try {
      const response = await getMe();
      set({ authUser: response.data, isLoading: false });

      // Set user identity for Amplitude analytics
      setUser(response.data.id);
      identifyUser({
        email: response.data.email,
        display_name: response.data.display_name || "Unknown",
      });

      return true;
    } catch (error) {
      logger.error("Access token validation failed:", error);
      if (refreshTokenValue) return await attemptTokenRefresh(set);
      return clearAuthenticationState(set);
    }
  }

  if (refreshTokenValue) return await attemptTokenRefresh(set);
  return clearAuthenticationState(set);
};

const attemptTokenRefresh = async (
  set: (partial: Partial<AuthState>) => void
): Promise<boolean> => {
  try {
    const response = await getMe();
    set({ authUser: response.data, isLoading: false });

    // Set user identity for Amplitude analytics
    setUser(response.data.id);
    identifyUser({
      email: response.data.email,
      display_name: response.data.display_name || "Unknown",
    });

    return true;
  } catch (error) {
    logger.error("Token refresh failed:", error);
    return clearAuthenticationState(set);
  }
};

const clearAuthenticationState = (
  set: (partial: Partial<AuthState>) => void
): boolean => {
  clearAuthTokens();
  set({ authUser: null, googleUser: null, isLoading: false });

  // Reset Amplitude user to anonymous
  setUser(null);

  return false;
};

const getErrorMessage = (error: unknown, defaultMessage: string): string =>
  error instanceof Error ? error.message : defaultMessage;

export const useAuthStore = create<AuthStore>()((set, get) => ({
  ...initialState,

  initializeAuth: async () => {
    const currentState = get();

    // Prevent duplicate initialization calls
    if (currentState.isInitializing) {
      return;
    }

    // Set initializing flag
    set({ isInitializing: true });

    try {
      await handleAuthentication(set);
    } finally {
      // Always clear the initializing flag
      set({ isInitializing: false });
    }
  },

  loginWithGoogle: async () => {
    try {
      set({ error: null });

      // Track login attempt
      authEvents.loginAttempt("google");

      const result = await signInWithPopup(auth, googleProvider);
      set({ googleUser: result.user });

      const idToken = await result.user.getIdToken();
      const refreshToken = result.user.refreshToken;

      const backendResponse = await authenticateWithGoogle({
        id_token: idToken,
        refresh_token: refreshToken,
      });

      setAccessToken(
        backendResponse.data.access_token,
        backendResponse.data.expires_in
      );
      if (backendResponse.data.refresh_token) {
        setRefreshToken(backendResponse.data.refresh_token);
      }

      set({
        googleUser: result.user,
        authUser: backendResponse.data.user,
      });

      // Track successful login and identify user
      authEvents.loginSuccess("google");
      setUser(backendResponse.data.user.id);
      identifyUser({
        email: backendResponse.data.user.email,
        display_name: backendResponse.data.user.display_name || "Unknown",
        signup_source: "google",
      });

      return backendResponse;
    } catch (error: unknown) {
      // Handle cancelled popup gracefully with comprehensive error detection
      if (
        error instanceof Error &&
        (error.message.includes("popup-closed-by-user") ||
          error.message.includes("auth/popup-closed-by-user") ||
          error.message.includes("auth/cancelled-popup-request") ||
          error.message.includes("auth/popup-blocked") ||
          error.message.includes("auth/popup-closed") ||
          (error.message.includes("popup") && error.message.includes("closed")))
      ) {
        logger.info("Google login popup closed by user");
        set({ error: null });
        return null;
      }

      logger.error("Google login error:", error);

      // Track login failure
      const errorMessage = getErrorMessage(
        error,
        "An error occurred during login"
      );
      authEvents.loginFailure("google", errorMessage);

      set({ error: errorMessage });
      throw error;
    }
  },

  loginWithGoogleRedirect: async () => {
    try {
      set({ isLoading: true, error: null });
      authEvents.loginAttempt("google");
      await signInWithRedirect(auth, googleProvider);
      // User will be redirected, no need to handle response here
    } catch (error: unknown) {
      logger.error("Google redirect login error:", error);
      const errorMessage = getErrorMessage(
        error,
        "An error occurred during login"
      );
      authEvents.loginFailure("google", errorMessage);
      set({ error: errorMessage, isLoading: false });
      throw error;
    }
  },

  handleRedirectResult: async () => {
    try {
      set({ isLoading: true, error: null });

      const result = await getRedirectResult(auth);
      if (!result) {
        // No redirect result, user didn't come from redirect flow
        set({ isLoading: false });
        return null;
      }

      set({ googleUser: result.user });

      const idToken = await result.user.getIdToken();
      const refreshToken = result.user.refreshToken;

      const backendResponse = await authenticateWithGoogle({
        id_token: idToken,
        refresh_token: refreshToken,
      });

      setAccessToken(
        backendResponse.data.access_token,
        backendResponse.data.expires_in
      );
      if (backendResponse.data.refresh_token) {
        setRefreshToken(backendResponse.data.refresh_token);
      }

      set({
        googleUser: result.user,
        authUser: backendResponse.data.user,
      });

      // Track successful login and identify user
      authEvents.loginSuccess("google");
      setUser(backendResponse.data.user.id);
      identifyUser({
        email: backendResponse.data.user.email,
        display_name: backendResponse.data.user.display_name || "Unknown",
        signup_source: "google",
      });

      set({ isLoading: false });
      return backendResponse;
    } catch (error: unknown) {
      logger.error("Handle redirect result error:", error);
      const errorMessage = getErrorMessage(
        error,
        "An error occurred during login"
      );
      authEvents.loginFailure("google", errorMessage);
      set({ error: errorMessage, isLoading: false });
      throw error;
    }
  },

  logout: async () => {
    try {
      set({ isLoading: true, error: null });
      clearAuthTokens();
      await signOut(auth);
      set({ googleUser: null, authUser: null, error: null });

      // Track logout and reset Amplitude user
      authEvents.logout();
      setUser(null);
    } catch (error: unknown) {
      logger.error("Logout error:", error);
      const errorMessage = getErrorMessage(
        error,
        "An error occurred during logout"
      );
      set({ error: errorMessage });
      throw error;
    } finally {
      set({ isLoading: false });
    }
  },

  refreshUserInfo: async () => {
    try {
      const response = await getMe();
      set({ authUser: response.data });
      return response.data;
    } catch (error: unknown) {
      logger.error("Failed to refresh user info:", error);
      throw error;
    }
  },

  checkAuthentication: async () => {
    return await handleAuthentication(set);
  },

  setGoogleUser: (user: FirebaseUser | null) =>
    set({ googleUser: user, error: null }),

  setAuthUser: (authUser: AuthUser | null) => set({ authUser, error: null }),

  setLoading: (isLoading: boolean) => set({ isLoading }),

  setError: (error: string | null) => set({ error, isLoading: false }),

  clearError: () => set({ error: null }),

  resetLoginState: () => set({ error: null, isLoading: false }),

  clearAuthState: () => {
    clearAuthTokens();
    set(initialState);
  },
}));
