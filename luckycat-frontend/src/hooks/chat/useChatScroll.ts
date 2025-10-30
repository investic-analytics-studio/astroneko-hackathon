import { useCallback, useEffect, useRef } from "react";
import type { FortuneMessage } from "@/hooks/fortune/useFortune";

export const useChatScroll = (messages: FortuneMessage[]) => {
  const chatRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = useCallback(() => {
    chatRef.current?.scrollTo({
      top: chatRef.current.scrollHeight,
      behavior: "smooth",
    });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  return chatRef;
};
