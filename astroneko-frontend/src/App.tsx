import "./App.css";
import { RouterProvider } from "@tanstack/react-router";
import { Suspense } from "react";
import { router } from "./routes/routes";
import { useEffect, useState } from "react";
import { logger } from "./lib/logger";
import { clearAllFortuneStates } from "./apis/fortune";
import { amplitudeInit, captureUTMParams } from "./lib/amplitude";
import { useAuthStore } from "./store/authStore";
import { googleAnalyticsInit } from "./config/ga-init";
import { useAuth } from "./hooks/common/useAuth";
import { AuthLoadingScreen } from "./components/common/AuthLoadingScreen";

function App() {
  const { isLoading: isAuthLoading } = useAuth();
  const [shouldShowLoading, setShouldShowLoading] = useState(true);

  useEffect(() => {
    if (!isAuthLoading) {
      const timer = setTimeout(() => {
        setShouldShowLoading(false);
      }, 100);

      return () => clearTimeout(timer);
    } else {
      setShouldShowLoading(true);
    }
  }, [isAuthLoading]);

  useEffect(() => {
    // Initialize Amplitude analytics
    amplitudeInit();

    // Initialize Google Analytics
    googleAnalyticsInit();

    // Capture UTM parameters from URL on first load
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.toString()) {
      captureUTMParams(urlParams);
    }

    // Handle Google Auth redirect result
    const handleAuthRedirect = async () => {
      try {
        const handleRedirectResult =
          useAuthStore.getState().handleRedirectResult;
        await handleRedirectResult();
      } catch (error) {
        logger.warn("Failed to handle redirect result:", error);
      }
    };

    const clearBackendSessions = async () => {
      try {
        await clearAllFortuneStates();
      } catch (error) {
        logger.warn("Failed to clear backend sessions on app init:", error);
      }
    };

    handleAuthRedirect();
    clearBackendSessions();
  }, []);

  // Show loading screen while auth is initializing or during delay
  if (isAuthLoading || shouldShowLoading) {
    return <AuthLoadingScreen />;
  }

  return (
    <Suspense
      fallback={
        <div className="min-h-screen relative bg-[image:var(--bg-display-2)] bg-cover bg-center overflow-hidden flex flex-col items-center justify-center">
          <div className="relative z-10 animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
        </div>
      }
    >
      <RouterProvider router={router} />
    </Suspense>
  );
}

export default App;
