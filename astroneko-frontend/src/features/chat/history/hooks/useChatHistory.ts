import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import { toast } from "sonner";
import { useTranslation } from "react-i18next";

import {
  getChatSessions,
  getChatHistory,
  deleteChatSession,
  type ChatHistoryParams,
  type ChatSession,
} from "@/apis/chat-history";
import { logger } from "@/lib/logger";

export const chatHistoryKeys = {
  all: ["chatHistory"] as const,
  sessions: () => [...chatHistoryKeys.all, "sessions"] as const,
  sessionsList: (params: ChatHistoryParams) =>
    [...chatHistoryKeys.sessions(), params] as const,
  session: (sessionId: string) =>
    [...chatHistoryKeys.all, "session", sessionId] as const,
  messages: (sessionId: string) =>
    [...chatHistoryKeys.session(sessionId), "messages"] as const,
} as const;

const TOAST_STYLES = {
  className: "custom-error-toast",
  descriptionClassName: "text-[#A1A1AA]",
  actionButtonStyle: { color: "#000000" },
} as const;

export const useChatSessions = (
  params: ChatHistoryParams = {},
  options: { enabled?: boolean } = {}
) => {
  return useQuery({
    queryKey: chatHistoryKeys.sessionsList(params),
    queryFn: async () => {
      const response = await getChatSessions(params);
      return response.data;
    },
    enabled: options.enabled !== false, // Default to true, but allow disabling
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes
    retry: (failureCount, error) => {
      // Don't retry on 4xx errors (client errors)
      if (error && typeof error === "object" && "status" in error) {
        const status = (error as any).status;
        if (status >= 400 && status < 500) return false;
      }
      return failureCount < 3;
    },
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
    refetchOnWindowFocus: false,
  });
};

export const useChatHistory = (sessionId: string, enabled: boolean = true) => {
  return useQuery({
    queryKey: chatHistoryKeys.messages(sessionId),
    queryFn: async () => {
      const response = await getChatHistory(sessionId);
      return response.data;
    },
    enabled: enabled && !!sessionId,
    staleTime: 0, // Cache for 5 minutes to prevent unnecessary refetches
    gcTime: 10 * 60 * 1000, // Keep in cache for 10 minutes
    retry: (failureCount, error) => {
      // Don't retry on 404 (session not found) or other 4xx errors
      if (error && typeof error === "object" && "status" in error) {
        const status = (error as any).status;
        if (status >= 400 && status < 500) return false;
      }
      return failureCount < 3;
    },
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    refetchInterval: false,
  });
};

// Deletes session with optimistic UI updates
export const useDeleteChatSession = () => {
  const queryClient = useQueryClient();
  const { t } = useTranslation();

  return useMutation({
    mutationFn: async (sessionId: string) => {
      const response = await deleteChatSession(sessionId);
      return response.data;
    },
    onMutate: async (sessionId: string) => {
      // Cancel any outgoing refetches
      await queryClient.cancelQueries({ queryKey: chatHistoryKeys.sessions() });

      // Snapshot the previous value
      const previousSessions = queryClient.getQueriesData({
        queryKey: chatHistoryKeys.sessions(),
      });

      // Optimistically update the cache
      queryClient.setQueriesData(
        { queryKey: chatHistoryKeys.sessions() },
        (old: { sessions: ChatSession[]; total: number } | undefined) => {
          if (!old) return old;
          return {
            ...old,
            sessions: old.sessions.filter(
              (session: ChatSession) => session.session_id !== sessionId
            ),
            total: old.total - 1,
          };
        }
      );

      return { previousSessions };
    },
    onError: (error, _sessionId, context) => {
      // Revert optimistic update on error
      if (context?.previousSessions) {
        context.previousSessions.forEach(([queryKey, data]) => {
          queryClient.setQueryData(queryKey, data);
        });
      }

      logger.error("Error deleting chat session:", error);
      toast.error(t("chat.error_delete_session"), TOAST_STYLES);
    },
    onSuccess: (_, sessionId) => {
      // Invalidate and refetch sessions
      queryClient.invalidateQueries({ queryKey: chatHistoryKeys.sessions() });

      // Remove the specific session from cache
      queryClient.removeQueries({
        queryKey: chatHistoryKeys.session(sessionId),
      });

      // Note: Toast will be handled by the calling component
      // to provide more specific messaging based on context
    },
  });
};

export const useRefreshChatSessions = () => {
  const queryClient = useQueryClient();

  return useCallback(() => {
    queryClient.invalidateQueries({ queryKey: chatHistoryKeys.sessions() });
  }, [queryClient]);
};

export const useRefetchChatHistory = () => {
  const queryClient = useQueryClient();

  return useCallback((sessionId: string) => {
    queryClient.invalidateQueries({
      queryKey: chatHistoryKeys.messages(sessionId),
    });
  }, [queryClient]);
};

export const usePrefetchChatHistory = () => {
  const queryClient = useQueryClient();

  return useCallback((sessionId: string) => {
    queryClient.prefetchQuery({
      queryKey: chatHistoryKeys.messages(sessionId),
      queryFn: async () => {
        const response = await getChatHistory(sessionId);
        return response.data;
      },
      staleTime: 2 * 60 * 1000,
    });
  }, [queryClient]);
};
