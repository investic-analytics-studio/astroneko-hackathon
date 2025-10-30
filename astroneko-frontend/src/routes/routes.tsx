import { createRouter, createRoute, redirect } from "@tanstack/react-router";
import { Route as rootRoute } from "./RouteLayout";
import { lazy } from "react";

// Lazy load pages for better initial bundle size
const LandingPage = lazy(() => import("../pages/LandingPage"));
const CategoryPage = lazy(() => import("../pages/CategoryPage"));
const PickCardTarot = lazy(() => import("../pages/PickCardPage"));
const ChatWithAI = lazy(() => import("../pages/ChatPage"));
const ActivateReferral = lazy(() => import("../pages/ActivateReferralPage"));
const GenesisPage = lazy(() => import("../pages/GenesisPage"));

// Language-based routes
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  beforeLoad: () => {
    // Redirect to default language using TanStack Router's redirect
    throw redirect({ to: "/$lng", params: { lng: "en" } });
  },
});

const languageRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/$lng",
  beforeLoad: ({ params }) => {
    // Validate language parameter
    const validLanguages = ["en", "th", "jp"];
    if (!validLanguages.includes(params.lng)) {
      throw redirect({ to: "/$lng", params: { lng: "en" } });
    }
  },
});

// Public routes
const landingRoute = createRoute({
  getParentRoute: () => languageRoute,
  path: "/",
  component: LandingPage,
});

// Public routes - removed protection to allow guest access
const categoryRoute = createRoute({
  getParentRoute: () => languageRoute,
  path: "/category",
  component: CategoryPage,
});

const pickCardRoute = createRoute({
  getParentRoute: () => languageRoute,
  path: "/pick-card",
  component: PickCardTarot,
});

const chatCategoryRoute = createRoute({
  getParentRoute: () => languageRoute,
  path: "/chat/$category",
  component: ChatWithAI,
});

const activateReferralRoute = createRoute({
  getParentRoute: () => languageRoute,
  path: "/activate-referral-code",
  component: ActivateReferral,
});

const genesisRoute = createRoute({
  getParentRoute: () => languageRoute,
  path: "/genesis",
  component: GenesisPage,
});

// Create route tree
const routeTree = rootRoute.addChildren([
  indexRoute,
  languageRoute.addChildren([
    landingRoute,
    categoryRoute,
    pickCardRoute,
    chatCategoryRoute,
    activateReferralRoute,
    genesisRoute,
  ]),
]);

// Create router
export const router = createRouter({
  routeTree,
});

// Register router for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
