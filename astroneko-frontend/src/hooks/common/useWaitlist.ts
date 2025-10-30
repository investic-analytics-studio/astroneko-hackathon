import { useWaitlistStore } from "@/store/waitlistStore";

export const useWaitlist = () => {
  const waitlistStore = useWaitlistStore();

  return {
    isOnWaitlist: waitlistStore.isOnWaitlist,
    isLoading: waitlistStore.isLoading,
    error: waitlistStore.error,
    waitlistData: waitlistStore.waitlistData,
    joinWaitlist: waitlistStore.joinWaitlist,
    checkWaitlistStatus: (email: string) =>
      waitlistStore.checkWaitlistStatus(email),
    clearError: waitlistStore.clearError,
  };
};
