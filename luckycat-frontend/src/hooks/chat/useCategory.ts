import { useMemo } from "react";
import { useLocation } from "@tanstack/react-router";

const DEFAULT_CATEGORY = "general";

const normalizeSearchParams = (search: unknown): URLSearchParams => {
  if (typeof search === "string") {
    return new URLSearchParams(search);
  }

  if (search && typeof search === "object") {
    const entries = Object.entries(search as Record<string, unknown>).reduce<
      Record<string, string>
    >((acc, [key, value]) => {
      if (value != null) {
        acc[key] = String(value);
      }
      return acc;
    }, {});

    return new URLSearchParams(entries);
  }

  return new URLSearchParams();
};

const resolveCategoryFromLocation = (
  pathname: string,
  search: unknown
): string => {
  const pathWithoutLang = pathname.replace(/^\/[a-z]{2}/, "") || "/";
  const isCategoryRoute =
    pathWithoutLang.startsWith("/chat/") && pathWithoutLang !== "/chat";

  if (isCategoryRoute) {
    const pathParts = pathWithoutLang.split("/");
    return (pathParts[2] || DEFAULT_CATEGORY).toLowerCase();
  }

  const searchParams = normalizeSearchParams(search);
  return (searchParams.get("category") || DEFAULT_CATEGORY).toLowerCase();
};

export const useCategory = () => {
  const location = useLocation();

  return useMemo(
    () => resolveCategoryFromLocation(location.pathname, location.search),
    [location.pathname, location.search]
  );
};
