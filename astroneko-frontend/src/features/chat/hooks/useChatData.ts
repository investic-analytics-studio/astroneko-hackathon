import { useMemo } from "react";
import { useChatHistory } from "../history";
import type { FortuneMessage } from "@/hooks/fortune/useFortune";
import { convertApiMessagesToFortuneMessages } from "../history/utils/chat-history-utils";

const deduplicateById = (messages: FortuneMessage[]): FortuneMessage[] => {
  const seen = new Set<string>();
  return messages.filter((msg) => {
    if (seen.has(msg.id)) return false;
    seen.add(msg.id);
    return true;
  });
};

export const useChatData = (
  sessionId: string | null,
  currentMessages: FortuneMessage[]
) => {
  const { data, isLoading } = useChatHistory(sessionId || "", !!sessionId);

  return useMemo(() => {
    const historyMessages = data?.messages
      ? convertApiMessagesToFortuneMessages(data.messages)
      : [];

    if (!sessionId) {
      return { messages: currentMessages, isLoading };
    }

    // When viewing history, combine history messages with new user messages
    // This ensures new user messages are displayed while preventing old session messages from persisting
    return {
      messages: deduplicateById([...historyMessages, ...currentMessages]),
      isLoading,
    };
  }, [data, currentMessages, sessionId, isLoading]);
};
