import { StateCreator } from "zustand";

export type ChatMode = "new" | "history";

export interface SessionSlice {
  currentSessionId: string | null;
  currentSessionName: string | null;
  chatMode: ChatMode;

  setCurrentSession: (
    sessionId: string | null,
    sessionName?: string | null
  ) => void;
  startNewChat: () => void;
  viewHistory: (sessionId: string, sessionName: string) => void;
}

export const createSessionSlice: StateCreator<SessionSlice> = (set) => ({
  currentSessionId: null,
  currentSessionName: null,
  chatMode: "new",

  setCurrentSession: (sessionId, sessionName) =>
    set({ currentSessionId: sessionId, currentSessionName: sessionName }),

  startNewChat: () =>
    set({
      currentSessionId: null,
      currentSessionName: null,
      chatMode: "new",
    }),

  viewHistory: (sessionId, sessionName) =>
    set({
      currentSessionId: sessionId,
      currentSessionName: sessionName,
      chatMode: "history",
    }),
});
