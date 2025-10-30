import { StateCreator } from "zustand";

export interface UISlice {
  isHistorySidebarOpen: boolean;
  searchQuery: string;

  setHistorySidebarOpen: (isOpen: boolean) => void;
  setSearchQuery: (query: string) => void;
  toggleHistorySidebar: () => void;
}

export const createUISlice: StateCreator<UISlice> = (set) => ({
  isHistorySidebarOpen: false,
  searchQuery: "",

  setHistorySidebarOpen: (isOpen) => set({ isHistorySidebarOpen: isOpen }),
  setSearchQuery: (query) => set({ searchQuery: query }),
  toggleHistorySidebar: () =>
    set((state) => ({ isHistorySidebarOpen: !state.isHistorySidebarOpen })),
});
