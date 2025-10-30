import { useEffect, useRef } from "react";
import { onAuthStateChanged, User as FirebaseUser } from "firebase/auth";
import { auth } from "@/config/firebase";
import { useAuthStore } from "@/store/authStore";
import { logger } from "@/lib/logger";

export const useAuth = () => {
  const authStore = useAuthStore();
  const hasInitialized = useRef(false);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(
      auth,
      async (googleUser: FirebaseUser | null) => {
        const store = useAuthStore.getState();
        store.setGoogleUser(googleUser);

        // Only initialize auth once on first load
        if (!hasInitialized.current && !store.authUser && !store.isInitializing) {
          logger.info("Initializing auth for the first time");
          hasInitialized.current = true;
          await store.initializeAuth();
        }
      }
    );

    return () => unsubscribe();
  }, []);

  return {
    googleUser: authStore.googleUser,
    authUser: authStore.authUser,
    isAuthenticated: !!authStore.authUser,
    isLoading: authStore.isLoading,
    error: authStore.error,
    loginWithGoogle: authStore.loginWithGoogle,
    logout: authStore.logout,
    refreshUserInfo: authStore.refreshUserInfo,
    checkAuthentication: authStore.checkAuthentication,
    setAuthUser: authStore.setAuthUser,
    clearError: authStore.clearError,
    resetLoginState: authStore.resetLoginState,
    clearAuthState: authStore.clearAuthState,
  };
};
