import { FortuneMessage } from "@/hooks/fortune/useFortune";
import { create } from "zustand";

export interface CategoryChatState {
  messages: FortuneMessage[];
  showQuestions: boolean;
  loading: boolean;
}

interface ChatActions {
  getCurrentCategoryState: (category: string) => CategoryChatState;
  setMessages: (category: string, messages: FortuneMessage[]) => void;
  addMessage: (category: string, message: FortuneMessage) => void;
  setShowQuestions: (category: string, show: boolean) => void;
  setLoading: (category: string, loading: boolean) => void;
  clearMessages: (category: string) => void;
  clearAllCategories: () => void;
  sendUserMessage: (category: string, message: string) => FortuneMessage;
  addAIMessage: (category: string, message: FortuneMessage) => void;
  resetCategory: (category: string) => void;
}

interface ChatState {
  categories: Record<string, CategoryChatState>;
}

type ChatStore = ChatState & ChatActions;

const defaultCategoryState: CategoryChatState = {
  messages: [],
  showQuestions: true,
  loading: false,
};

const initialState: ChatState = {
  categories: {},
};

export const useChatStore = create<ChatStore>()((set, get) => ({
  ...initialState,

  getCurrentCategoryState: (category: string) => {
    return get().categories[category] || defaultCategoryState;
  },

  setMessages: (category: string, messages: FortuneMessage[]) =>
    set((state) => ({
      categories: {
        ...state.categories,
        [category]: {
          ...state.categories[category],
          messages,
        },
      },
    })),

  addMessage: (category: string, message: FortuneMessage) =>
    set((state) => {
      const current = state.categories[category] || defaultCategoryState;
      return {
        categories: {
          ...state.categories,
          [category]: {
            ...current,
            messages: [...current.messages, message],
          },
        },
      };
    }),

  setShowQuestions: (category: string, show: boolean) =>
    set((state) => ({
      categories: {
        ...state.categories,
        [category]: {
          ...(state.categories[category] || defaultCategoryState),
          showQuestions: show,
        },
      },
    })),

  setLoading: (category: string, loading: boolean) =>
    set((state) => ({
      categories: {
        ...state.categories,
        [category]: {
          ...(state.categories[category] || defaultCategoryState),
          loading,
        },
      },
    })),

  clearMessages: (category: string) =>
    set((state) => ({
      categories: {
        ...state.categories,
        [category]: defaultCategoryState,
      },
    })),

  clearAllCategories: () => set(initialState),

  sendUserMessage: (category: string, message: string): FortuneMessage => {
    const userMessage: FortuneMessage = {
      role: "user",
      message: message,
      id: `user-${Date.now()}`,
    };

    get().addMessage(category, userMessage);
    return userMessage;
  },

  addAIMessage: (category: string, message: FortuneMessage) => {
    get().addMessage(category, message);
  },

  resetCategory: (category: string) => {
    set((state) => ({
      categories: {
        ...state.categories,
        [category]: defaultCategoryState,
      },
    }));
  },
}));
