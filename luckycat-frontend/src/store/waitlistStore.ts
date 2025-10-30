import { create } from "zustand";
import { logger } from "../lib/logger";
import {
  addToWaitingList,
  checkWaitingList,
  type WaitingListData,
  type WaitingListCheckData,
} from "../apis/waitinglist";

interface WaitlistData extends WaitingListData {
  is_in_waiting_list?: boolean;
  [key: string]: unknown;
}

interface WaitlistState {
  isOnWaitlist: boolean;
  isLoading: boolean;
  error: string | null;
  waitlistData: WaitlistData | null;
}

interface WaitlistActions {
  joinWaitlist: (email: string) => Promise<void>;
  checkWaitlistStatus: (email: string) => Promise<WaitlistData | null>;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  reset: () => void;
}

type WaitlistStore = WaitlistState & WaitlistActions;

const initialState: WaitlistState = {
  isOnWaitlist: false,
  isLoading: false,
  error: null,
  waitlistData: null,
};

const getErrorMessage = (error: unknown, defaultMessage: string): string =>
  error instanceof Error ? error.message : defaultMessage;

const handleApiResponse = (data: WaitingListData): WaitlistData =>
  data as WaitlistData;

export const useWaitlistStore = create<WaitlistStore>((set) => ({
  ...initialState,

  joinWaitlist: async (email: string) => {
    try {
      set({ isLoading: true, error: null });

      const response = await addToWaitingList({ email });
      const waitlistData = handleApiResponse(response.data);

      set({
        isOnWaitlist: true,
        waitlistData,
        isLoading: false,
      });
    } catch (error: unknown) {
      const errorMessage = getErrorMessage(error, "Failed to join waitlist");
      logger.error("Failed to join waitlist:", error);

      set({
        error: errorMessage,
        isLoading: false,
      });

      throw error;
    }
  },

  checkWaitlistStatus: async (email: string) => {
    try {
      set({ isLoading: true, error: null });

      const response = await checkWaitingList({ email });
      const checkData = response.data as WaitingListCheckData;

      const isInWaitingList = checkData.message === "User is in waiting list";

      const waitlistData: WaitlistData = {
        id: "",
        email: email,
        created_at: "",
        updated_at: "",
        is_in_waiting_list: isInWaitingList,
      };

      set({
        isOnWaitlist: isInWaitingList,
        waitlistData,
        isLoading: false,
      });

      return waitlistData;
    } catch (error: unknown) {
      const errorMessage = getErrorMessage(
        error,
        "Failed to check waitlist status"
      );
      logger.error("Failed to check waitlist status:", error);

      set({
        error: errorMessage,
        isLoading: false,
      });

      throw error;
    }
  },

  setLoading: (isLoading: boolean) => set({ isLoading }),

  setError: (error: string | null) => set({ error, isLoading: false }),

  clearError: () => set({ error: null }),

  reset: () => set(initialState),
}));
