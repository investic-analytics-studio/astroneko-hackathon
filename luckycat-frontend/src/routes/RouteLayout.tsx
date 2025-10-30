import { createRootRoute, Outlet, useParams } from "@tanstack/react-router";
import { MainLayout } from "../layouts/MainLayout";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { usePageView } from "../hooks/common/usePageView";

const LanguageWrapper = () => {
  const params = useParams({ from: "/$lng" }) as { lng?: string };
  const { i18n } = useTranslation();

  // Track page views automatically
  usePageView();

  useEffect(() => {
    if (params.lng && ["en", "th", "jp"].includes(params.lng)) {
      i18n.changeLanguage(params.lng);
      localStorage.setItem("language", params.lng);
    }
  }, [params.lng, i18n]);

  return (
    <MainLayout>
      <Outlet />
    </MainLayout>
  );
};

export const Route = createRootRoute({
  component: LanguageWrapper,
});
