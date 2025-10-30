import { useChatHistoryStore } from "./chatHistoryStore";

type StoreState = ReturnType<typeof useChatHistoryStore.getState>;

// Backward compatible selectors
export const chatHistorySelectors = {
  currentSessionId: (state: StoreState) => state.currentSessionId,
  currentSessionName: (state: StoreState) => state.currentSessionName,
  chatMode: (state: StoreState) => state.chatMode,
  isHistorySidebarOpen: (state: StoreState) => state.isHistorySidebarOpen,
  searchQuery: (state: StoreState) => state.searchQuery,
  sessions: (state: StoreState) => state.sessions,
  currentMessages: (state: StoreState) => state.currentMessages,
};
