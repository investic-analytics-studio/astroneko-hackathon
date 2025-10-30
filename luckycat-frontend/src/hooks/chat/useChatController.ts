import { useState } from "react";
import { useImagePreloader } from "@/hooks";
import { useCategory } from "@/hooks";
import { useChatMessages } from "@/hooks/chat/useChatMessages";
import { useChatScroll } from "@/hooks/chat/useChatScroll";
import { useChatSuggestions } from "@/hooks/chat/useChatSuggestions";
import { useChatHistoryStore } from "@/features/chat/history";
import type { CategoryChatState } from "@/store/chatStore";
import { track } from "@/lib/amplitude";

const PRELOAD_IMAGES = ["icons/cat-foot.webp", "cards/card-back.webp"];

export interface UseChatControllerReturn {
  chatRef: React.RefObject<HTMLDivElement | null>;
  input: string;
  setInput: (value: string) => void;
  category: string;
  messages: CategoryChatState["messages"];
  loading: boolean;
  showQuestions: boolean;
  categoryQuestions: string[];
  toggleQuestions: () => void;
  handleSelectSuggestion: (message: string) => void;
  handleSendMessage: (message: string) => Promise<void>;
  resetChat: () => Promise<void>;
}

export const useChatController = (): UseChatControllerReturn => {
  const [input, setInput] = useState("");

  const category = useCategory();
  const { messages, loading, send, reset } = useChatMessages(category);
  const chatRef = useChatScroll(messages);
  const { questions, showQuestions, toggle, hide } =
    useChatSuggestions(category);
  const startNewChat = useChatHistoryStore((state) => state.startNewChat);

  useImagePreloader({ images: PRELOAD_IMAGES, priority: false });

  const handleSelectSuggestion = (message: string) => {
    // Track suggestion selection event
    track('suggestion selected', {
      category,
      suggestion_en: message, // English question text
      suggestion_length: message.length,
      question_type: 'suggested',
    });

    // Track specific question event with question text as event name
    const cleanQuestionText = message
      .trim()
      .replace(/[?!.,]/g, '') // Remove punctuation
      .replace(/\s+/g, ' ') // Normalize spaces
      .substring(0, 100); // Limit length

    const questionEventName = `Question Asked - ${cleanQuestionText}`;
    track(questionEventName, {
      category,
      question_text: message,
      question_length: message.length,
      question_type: 'suggested',
      selection_method: 'suggestion_click',
    });

    setInput(message);
    hide();
  };

  const handleSendMessage = async (message: string) => {
    // Track general message sent event
    track('message sent', {
      category,
      message_length: message.length,
      is_first_message: messages.length === 0,
    });

    setInput(""); // Clear input immediately
    hide(); // Hide chat suggestions when sending a message
    await send(message);
  };

  const resetChat = async () => {
    track('chat reset', {
      category,
      message_count: messages.length,
    });
    await reset();
    setInput("");
    startNewChat();
    chatRef.current?.scrollTo({ top: 0, behavior: "smooth" });
  };

  return {
    chatRef,
    input,
    setInput,
    category,
    messages,
    loading,
    showQuestions,
    categoryQuestions: questions,
    toggleQuestions: toggle,
    handleSelectSuggestion,
    handleSendMessage,
    resetChat,
  };
};
