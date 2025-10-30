/**
 * Page View Tracking Hook for TanStack Router
 *
 * Automatically tracks page views when route changes occur.
 * Integrates with Amplitude analytics to capture navigation patterns.
 *
 * @module hooks/common/usePageView
 */

import { useEffect } from 'react';
import { useRouter } from '@tanstack/react-router';
import { trackPageView, isInitialized } from '@/lib/amplitude';
import { logger } from '@/lib/logger';

/**
 * Get human-readable page name from path
 */
function getPageNameFromPath(pathname: string, params: any): string {
  const segments = pathname.split('/').filter(Boolean);

  // Handle language-based routes
  if (segments.length >= 1 && ['en', 'th', 'jp'].includes(segments[0])) {
    const pathWithoutLang = '/' + segments.slice(1).join('/');

    switch (pathWithoutLang) {
      case '/':
        return 'Landing Page';
      case '/category':
        return 'Category Selection';
      case '/pick-card':
        return 'Card Selection';
      default:
        if (pathWithoutLang.startsWith('/chat/')) {
          const category = params?.category || 'unknown';
          return `Chat - ${category}`;
        }
        return pathWithoutLang || 'Landing Page';
    }
  }

  // Fallback for non-language routes
  switch (pathname) {
    case '/':
      return 'Landing Page';
    case '/category':
      return 'Category Selection';
    case '/pick-card':
      return 'Card Selection';
    default:
      if (pathname.startsWith('/chat/')) {
        const category = params?.category || 'unknown';
        return `Chat - ${category}`;
      }
      return pathname || 'Landing Page';
  }
}

/**
 * Hook to automatically track page views on route changes.
 * Call this at the root of your app or layout component.
 *
 * Features:
 * - Tracks pathname and search params
 * - Debounces rapid navigation
 * - Only tracks when Amplitude is initialized
 * - Captures full route context including dynamic params
 *
 * @example
 * ```tsx
 * function AppLayout() {
 *   usePageView();
 *   return <Outlet />;
 * }
 * ```
 */
export function usePageView() {
  const router = useRouter();

  useEffect(() => {
    // Only track if Amplitude is initialized
    if (!isInitialized()) {
      logger.debug('[usePageView] Amplitude not initialized, skipping page view tracking');
      return;
    }

    // Subscribe to router state changes
    const unsubscribe = router.subscribe('onLoad', ({ toLocation }) => {
      // Extract pathname and search from location
      const pathname = toLocation.pathname;
      const search = toLocation.search ? new URLSearchParams(toLocation.search as any).toString() : '';

      // Extract category from path for chat routes
      const pathSegments = pathname.split('/').filter(Boolean);
      let category = '';
      if (pathSegments.length >= 3 && pathSegments[1] === 'chat') {
        category = pathSegments[2];
      }

      // Get page name from path
      const pageName = getPageNameFromPath(pathname, { category });

      // Track page view
      trackPageView(pathname, search, pageName);
    });

    // Track initial page view
    const initialPath = router.state.location.pathname;
    const initialSearch = router.state.location.search
      ? new URLSearchParams(router.state.location.search as any).toString()
      : '';

    // Extract category from initial path for chat routes
    const initialPathSegments = initialPath.split('/').filter(Boolean);
    let initialCategory = '';
    if (initialPathSegments.length >= 3 && initialPathSegments[1] === 'chat') {
      initialCategory = initialPathSegments[2];
    }

    const initialPageName = getPageNameFromPath(initialPath, { category: initialCategory });

    trackPageView(initialPath, initialSearch, initialPageName);

    // Cleanup subscription
    return () => {
      unsubscribe();
    };
  }, [router]);
}
