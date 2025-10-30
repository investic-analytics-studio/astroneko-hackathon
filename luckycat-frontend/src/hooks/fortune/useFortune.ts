import { useCallback } from "react";
import { useTranslation } from "react-i18next";
import {
  clearFortuneState,
  getFortuneReply,
  type FortuneResponse,
  type FortuneCategory,
} from "@/apis/fortune";
import { FORTUNE_CARDS } from "@/constants/fortune";
import { handleFortuneError } from "@/lib/error-handlers";
import { logger } from "@/lib/logger";

export interface FortuneMessage {
  role: "user" | "ai";
  message: string;
  id: string;
  card?: string;
  meaning?: string;
}

export interface UseFortuneReturn {
  sendMessage: (
    message: string,
    category: string,
    chatHistorySessionId?: string | null
  ) => Promise<FortuneMessage | null>;
  clearFortuneSession: (category: string) => Promise<void>;
  clearAllFortuneSessions: () => Promise<void>;
}

const DEFAULT_ERROR_MESSAGE =
  "Sorry, I encountered an error. Please try again later.";

const getCardDisplayName = (cardName?: string): string | undefined => {
  if (!cardName) return undefined;
  return FORTUNE_CARDS.find((card) => card.name === cardName)?.display_name;
};

const createFortuneMessage = (
  role: "user" | "ai",
  message: string,
  options?: { card?: string; meaning?: string }
): FortuneMessage => ({
  role,
  message,
  id: `${role}-${Date.now()}`,
  card: options?.card,
  meaning: options?.meaning,
});

export const useFortune = (): UseFortuneReturn => {
  const { t } = useTranslation();

  const sendMessage = useCallback(
    async (
      message: string,
      category: string,
      chatHistorySessionId?: string | null
    ): Promise<FortuneMessage | null> => {
      if (message.trim() === "") return null;

      try {
        const response: FortuneResponse = await getFortuneReply(
          message,
          category,
          chatHistorySessionId
        );

        return createFortuneMessage("ai", response.message, {
          card: getCardDisplayName(response.card),
          meaning: response.meaning,
        });
      } catch (error: unknown) {
        const errorMessage = handleFortuneError(
          error,
          t,
          DEFAULT_ERROR_MESSAGE
        );
        return createFortuneMessage("ai", errorMessage);
      }
    },
    [t]
  );

  const clearFortuneSession = useCallback(
    async (category: string): Promise<void> => {
      try {
        await clearFortuneState(category as FortuneCategory);
      } catch (error) {
        logger.error(
          `Error clearing fortune state for category ${category}:`,
          error
        );
      }
    },
    []
  );

  const clearAllFortuneSessions = useCallback(async (): Promise<void> => {
    const categories: FortuneCategory[] = ["general", "crypto", "lover"];

    const clearPromises = categories.map(async (category) => {
      try {
        await clearFortuneState(category);
      } catch (error) {
        logger.error(
          `Failed to clear fortune state for category ${category}:`,
          error
        );
      }
    });

    await Promise.all(clearPromises);
  }, []);

  return {
    sendMessage,
    clearFortuneSession,
    clearAllFortuneSessions,
  };
};
