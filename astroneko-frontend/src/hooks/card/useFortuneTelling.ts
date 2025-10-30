import { useState, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { clearFortuneState, getFortuneReply } from "@/apis/fortune";
import { handleFortuneError } from "@/lib/error-handlers";
import { logger } from "@/lib/logger";

const DEFAULT_ERROR_MESSAGE =
  "The cards are mysterious at this moment. Try again later.";

export const useFortuneTelling = () => {
  const { t } = useTranslation();
  const [fortuneReply, setFortuneReply] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);

  const getFortune = useCallback(
    async (prompt: string) => {
      setIsLoading(true);

      try {
        const reply = await getFortuneReply(prompt);
        setFortuneReply(reply.message);
        await clearFortuneState();
      } catch (error) {
        logger.error("Error getting fortune:", error);
        const errorMessage = handleFortuneError(
          error,
          t,
          DEFAULT_ERROR_MESSAGE
        );
        setFortuneReply(errorMessage);
      } finally {
        setIsLoading(false);
      }
    },
    [t]
  );

  const clearReply = useCallback(() => {
    setFortuneReply("");
  }, []);

  return {
    fortuneReply,
    isLoading,
    getFortune,
    clearReply,
  };
};
