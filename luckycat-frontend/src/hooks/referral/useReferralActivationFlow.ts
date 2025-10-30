import { useEffect, useState, useRef } from "react";
import { useNavigate, useParams } from "@tanstack/react-router";
import { useAuthStore } from "@/store/authStore";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";
import { logger } from "@/lib/logger";
import { detectWebView } from "@/lib/webviewDetection";
import { useReferralActivation } from "./useReferralActivation";

/**
 * Orchestrates the referral activation flow
 * Handles: URL parsing, auth checks, modal management, navigation
 */
export function useReferralActivationFlow() {
  const navigate = useNavigate();
  const { lng } = useParams({ strict: false });
  const { t } = useTranslation();
  const { authUser: user } = useAuthStore();
  const { activate, isLoading, isSuccess, isError } = useReferralActivation();

  const [referralCode, setReferralCode] = useState("");
  const [loginModalOpen, setLoginModalOpen] = useState(false);
  const [isNavigating, setIsNavigating] = useState(false);
  const hasActivatedRef = useRef(false);

  // Extract referral code from URL on mount
  useEffect(() => {
    const webViewInfo = detectWebView();
    if (webViewInfo.isWebView) {
      logger.warn("Referral activation blocked in in-app browser");
      return;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get("code");

    if (!code) {
      logger.warn("No referral code found in URL");
      navigate({ to: "/$lng", params: { lng: lng || "en" } });
      return;
    }

    logger.info("Referral code found:", code);
    setReferralCode(code);
  }, [navigate, lng]);

  // Show login modal if user is not authenticated
  useEffect(() => {
    if (referralCode && !user && !loginModalOpen && !isNavigating) {
      const timer = setTimeout(() => {
        if (!user && !isNavigating) {
          logger.info("User not logged in, showing login modal");
          setLoginModalOpen(true);
        }
      }, 200);
      return () => clearTimeout(timer);
    }
  }, [referralCode, user, loginModalOpen, isNavigating]);

  // Auto-close login modal when user logs in
  useEffect(() => {
    if (user && loginModalOpen) {
      logger.info("User logged in, closing login modal");
      const timer = setTimeout(() => {
        setLoginModalOpen(false);
      }, 300);
      return () => clearTimeout(timer);
    }
  }, [user, loginModalOpen]);

  // Handle activation when user is logged in
  useEffect(() => {
    if (!user || !referralCode || hasActivatedRef.current) {
      return;
    }

    // Check if already activated
    if (user.is_activated_referral) {
      logger.info("User already has activated referral");
      toast.info(
        t("referral.alreadyActivated", "Referral code already activated")
      );
      navigate({ to: "/$lng", params: { lng: lng || "en" }, search: {} });
      return;
    }

    // Start activation
    logger.info("Starting referral activation for user:", user.id);
    hasActivatedRef.current = true;
    activate({ referral_code: referralCode });
  }, [user, referralCode]);

  // Handle success/error navigation
  useEffect(() => {
    if (isSuccess) {
      const timer = setTimeout(() => {
        navigate({ to: "/$lng", params: { lng: lng || "en" }, search: {} });
      }, 2000);
      return () => clearTimeout(timer);
    }

    if (isError) {
      const timer = setTimeout(() => {
        navigate({ to: "/$lng", params: { lng: lng || "en" }, search: {} });
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [isSuccess, isError]);

  const handleLoginModalClose = (open: boolean) => {
    if (!open && !user) {
      logger.info("User closed login modal without logging in");
      setIsNavigating(true);
      // Add smooth transition delay before navigation
      setTimeout(() => {
        navigate({ to: "/$lng", params: { lng: lng || "en" } });
      }, 300);
    }
    setLoginModalOpen(open);
  };

  const handleBackToHome = () => {
    navigate({ to: "/$lng", params: { lng: lng || "en" } });
  };

  return {
    loginModalOpen,
    referralCode,
    isLoading,
    isSuccess,
    isError,
    isNavigating,
    handleLoginModalClose,
    handleBackToHome,
  };
}
