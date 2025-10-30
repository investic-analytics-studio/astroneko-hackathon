import { create } from "zustand";
import { type ReferralCode } from "../apis/auth";

// State interface
interface ReferralState {
  codes: ReferralCode[];
  isLoading: boolean;
  error: string | null;
}

// Actions interface
interface ReferralActions {
  setCodes: (codes: ReferralCode[]) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  reset: () => void;
}

type ReferralStore = ReferralState & ReferralActions;

// Initial state constant
const initialState: ReferralState = {
  codes: [],
  isLoading: false,
  error: null,
};

export const useReferralStore = create<ReferralStore>((set) => ({
  // State
  ...initialState,

  // Actions
  setCodes: (codes: ReferralCode[]) => set({ codes, error: null }),

  setLoading: (isLoading: boolean) => set({ isLoading }),

  setError: (error: string | null) => set({ error, isLoading: false }),

  clearError: () => set({ error: null }),

  reset: () => set(initialState),
}));
