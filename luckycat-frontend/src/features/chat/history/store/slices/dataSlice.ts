import { StateCreator } from "zustand";
import type { ChatSession, ChatMessage } from "@/apis/chat-history";

export interface DataSlice {
  sessions: ChatSession[];
  currentMessages: ChatMessage[];

  setSessions: (sessions: ChatSession[]) => void;
  addSession: (session: ChatSession) => void;
  removeSession: (sessionId: string) => void;
  updateSession: (sessionId: string, updates: Partial<ChatSession>) => void;
  setCurrentMessages: (messages: ChatMessage[]) => void;
  clearCurrentMessages: () => void;
}

export const createDataSlice: StateCreator<DataSlice> = (set) => ({
  sessions: [],
  currentMessages: [],

  setSessions: (sessions) => set({ sessions }),

  addSession: (session) =>
    set((state) => ({
      sessions: [session, ...state.sessions],
    })),

  removeSession: (sessionId) =>
    set((state) => ({
      sessions: state.sessions.filter((s) => s.session_id !== sessionId),
    })),

  updateSession: (sessionId, updates) =>
    set((state) => ({
      sessions: state.sessions.map((s) =>
        s.session_id === sessionId ? { ...s, ...updates } : s
      ),
    })),

  setCurrentMessages: (messages) => set({ currentMessages: messages }),
  clearCurrentMessages: () => set({ currentMessages: [] }),
});
