import { create } from "zustand";
import { devtools } from "zustand/middleware";
import { createSessionSlice, SessionSlice } from "./slices/sessionSlice";
import { createUISlice, UISlice } from "./slices/uiSlice";
import { createDataSlice, DataSlice } from "./slices/dataSlice";

export type { ChatMode } from "./slices/sessionSlice";

type ChatHistoryStore = SessionSlice & UISlice & DataSlice;

export const useChatHistoryStore = create<ChatHistoryStore>()(
  devtools(
    (...a) => ({
      ...createSessionSlice(...a),
      ...createUISlice(...a),
      ...createDataSlice(...a),
    }),
    { name: "chat-history" }
  )
);

// Simple selectors (use with shallow comparison in components)
export const selectSession = (state: ChatHistoryStore) => ({
  currentSessionId: state.currentSessionId,
  currentSessionName: state.currentSessionName,
  chatMode: state.chatMode,
});

export const selectUI = (state: ChatHistoryStore) => ({
  isHistorySidebarOpen: state.isHistorySidebarOpen,
  searchQuery: state.searchQuery,
});

export const selectSessions = (state: ChatHistoryStore) => state.sessions;
export const selectCurrentMessages = (state: ChatHistoryStore) =>
  state.currentMessages;
