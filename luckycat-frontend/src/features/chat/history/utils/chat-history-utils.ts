import type { ChatMessage as ApiChatMessage } from "@/apis/chat-history";
import type { FortuneMessage } from "@/hooks/fortune/useFortune";

export const convertApiMessageToFortuneMessage = (
  apiMessage: ApiChatMessage
): FortuneMessage => ({
  role: apiMessage.role,
  message: apiMessage.message,
  id: apiMessage.id,
  card: apiMessage.card,
  meaning: apiMessage.meaning,
});

export const convertApiMessagesToFortuneMessages = (
  apiMessages: ApiChatMessage[]
): FortuneMessage[] => {
  return apiMessages.map(convertApiMessageToFortuneMessage);
};

/**
 * Check if a message has card data
 * @param message - Message to check
 * @returns True if message has card data
 */
export const hasCardData = (
  message: FortuneMessage | ApiChatMessage
): boolean => {
  return !!(message.card && message.meaning);
};

/**
 * Format message timestamp for display
 * @param timestamp - ISO timestamp string
 * @returns Formatted timestamp string
 */
export const formatMessageTimestamp = (timestamp: string): string => {
  const date = new Date(timestamp);
  return date.toLocaleString();
};

/**
 * Get message age in human readable format
 * @param timestamp - ISO timestamp string
 * @returns Human readable age string
 */
export const getMessageAge = (timestamp: string): string => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffInMs = now.getTime() - date.getTime();
  const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
  const diffInHours = Math.floor(diffInMinutes / 60);
  const diffInDays = Math.floor(diffInHours / 24);

  if (diffInMinutes < 1) return "Just now";
  if (diffInMinutes < 60) return `${diffInMinutes}m ago`;
  if (diffInHours < 24) return `${diffInHours}h ago`;
  if (diffInDays < 7) return `${diffInDays}d ago`;

  return date.toLocaleDateString();
};
