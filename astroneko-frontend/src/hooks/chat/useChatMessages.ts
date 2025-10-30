import { useCallback, useRef, useEffect } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useChatStore } from "@/store/chatStore";
import { useFortune } from "@/hooks/fortune/useFortune";
import {
  useChatHistoryStore,
  chatHistorySelectors,
  chatHistoryKeys,
} from "@/features/chat/history";
import { logger } from "@/lib/logger";

export const useChatMessages = (category: string) => {
  const sessionRefreshTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const queryClient = useQueryClient();

  const {
    getCurrentCategoryState,
    setLoading,
    sendUserMessage,
    addAIMessage,
    resetCategory,
  } = useChatStore();

  const { messages, loading } = getCurrentCategoryState(category);
  const { sendMessage: sendFortuneMessage, clearFortuneSession } = useFortune();

  const currentSessionId = useChatHistoryStore(
    chatHistorySelectors.currentSessionId
  );

  const send = useCallback(
    async (message: string) => {
      const trimmedMessage = message.trim();

      if (trimmedMessage === "") {
        logger.warn("Cannot send empty message");
        return;
      }

      if (trimmedMessage.length > 2000) {
        logger.error("Message too long (max 2000 characters)");
        return;
      }

      sendUserMessage(category, trimmedMessage);
      setLoading(category, true);

      try {
        const sessionIdToUse = currentSessionId || null;
        const aiResponse = await sendFortuneMessage(
          trimmedMessage,
          category,
          sessionIdToUse
        );

        if (aiResponse) {
          addAIMessage(category, aiResponse);

          if (sessionRefreshTimeoutRef.current) {
            clearTimeout(sessionRefreshTimeoutRef.current);
          }

          sessionRefreshTimeoutRef.current = setTimeout(() => {
            queryClient.invalidateQueries({
              queryKey: chatHistoryKeys.sessions(),
            });
            sessionRefreshTimeoutRef.current = null;
          }, 2000);
        }
      } catch (error) {
        logger.error("Error in sendMessage:", error);
      } finally {
        setLoading(category, false);
      }
    },
    [
      category,
      currentSessionId,
      sendUserMessage,
      setLoading,
      sendFortuneMessage,
      addAIMessage,
      queryClient,
    ]
  );

  const reset = useCallback(async () => {
    try {
      await clearFortuneSession(category);
    } catch (error) {
      logger.error("Error resetting chat:", error);
    } finally {
      resetCategory(category);
    }
  }, [category, clearFortuneSession, resetCategory]);

  useEffect(() => {
    return () => {
      if (sessionRefreshTimeoutRef.current) {
        clearTimeout(sessionRefreshTimeoutRef.current);
      }
    };
  }, []);

  return { messages, loading, send, reset };
};
