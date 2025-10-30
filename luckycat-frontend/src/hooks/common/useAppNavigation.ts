import { useNavigate } from "@tanstack/react-router";
import { useLanguage } from "./useLanguage";

export const useAppNavigation = () => {
  const navigate = useNavigate();
  const { currentLanguage } = useLanguage();

  return {
    goToHome: () => navigate({ to: "/$lng", params: { lng: currentLanguage } }),
    goToCategory: () =>
      navigate({ to: "/$lng/category", params: { lng: currentLanguage } }),
    goToPickCard: () =>
      navigate({ to: "/$lng/pick-card", params: { lng: currentLanguage } }),
  };
};
