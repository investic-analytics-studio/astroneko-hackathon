import { useState, useEffect } from "react";

interface DeviceDetection {
  isMobile: boolean;
  isTablet: boolean;
  isDesktop: boolean;
}

// Breakpoints for device detection
const MOBILE_BREAKPOINT = 768;
const TABLET_BREAKPOINT = 1024;

/**
 * Comprehensive device detection hook that combines mobile detection, responsive breakpoints, and user agent detection
 * Uses matchMedia API for better performance and accuracy
 */
export function useDeviceDetection(): DeviceDetection {
  const [deviceType, setDeviceType] = useState<DeviceDetection>(() => {
    // Initial server-side safe detection
    if (typeof window === "undefined") {
      return {
        isMobile: false,
        isTablet: false,
        isDesktop: true,
      };
    }

    return {
      isMobile: window.innerWidth < MOBILE_BREAKPOINT,
      isTablet: window.innerWidth >= MOBILE_BREAKPOINT && window.innerWidth <= TABLET_BREAKPOINT,
      isDesktop: window.innerWidth > TABLET_BREAKPOINT,
    };
  });

  useEffect(() => {
    if (typeof window === "undefined") return;

    const mobileMql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`);
    const tabletMql = window.matchMedia(`(min-width: ${MOBILE_BREAKPOINT}px) and (max-width: ${TABLET_BREAKPOINT}px)`);

    const checkDevice = () => {
      const width = window.innerWidth;

      // Check responsive breakpoints
      const isMobileByWidth = width < MOBILE_BREAKPOINT;
      const isTabletByWidth = width >= MOBILE_BREAKPOINT && width <= TABLET_BREAKPOINT;
      const isDesktopByWidth = width > TABLET_BREAKPOINT;

      // Check user agent for mobile devices (helps with tablets that might report as mobile)
      const isMobileByUA = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
        navigator.userAgent
      );

      // Combine both detection methods for better accuracy
      const isMobile = isMobileByWidth || (isMobileByUA && width <= TABLET_BREAKPOINT);
      const isTablet = isTabletByWidth && !isMobile;
      const isDesktop = isDesktopByWidth && !isMobile && !isTablet;

      setDeviceType({
        isMobile,
        isTablet,
        isDesktop,
      });
    };

    // Check initially
    checkDevice();

    // Add event listeners for better performance using matchMedia
    const handleMobileChange = () => checkDevice();
    const handleTabletChange = () => checkDevice();

    mobileMql.addEventListener("change", handleMobileChange);
    tabletMql.addEventListener("change", handleTabletChange);

    // Cleanup
    return () => {
      mobileMql.removeEventListener("change", handleMobileChange);
      tabletMql.removeEventListener("change", handleTabletChange);
    };
  }, []);

  return deviceType;
}

/**
 * Simple mobile detection hook for backward compatibility
 * @returns boolean indicating if the current device is mobile
 */
export function useMobile(): boolean {
  const { isMobile } = useDeviceDetection();
  return isMobile;
}

/**
 * Alias for useMobile for backward compatibility
 */
export const useIsMobile = useMobile;

/**
 * Mobile detection hook with UA string for more accurate detection
 * @returns boolean indicating if the current device is mobile
 */
export function useMobileDetection(): boolean {
  return useMobile();
}
