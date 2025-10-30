import { API_ENDPOINTS } from "@/config/api";
import type { ApiResponse } from "@/types/api";
import axios from "../config/axios";

export interface ChatSession {
  session_id: string;
  history_name: string;
  created_at: string;
  updated_at: string;
}

export interface ChatMessage {
  id: string;
  message: string;
  role: "user" | "ai";
  used_tokens: number;
  created_at: string;
  card?: string;
  meaning?: string;
}

export interface ChatHistoryResponse {
  session_id: string;
  history_name: string;
  messages: ChatMessage[];
  total: number;
}

export interface ChatSessionsResponse {
  sessions: ChatSession[];
  total: number;
}

// Note: search param is deprecated - filtering done client-side for better UX
export interface ChatHistoryParams {
  sort_by?: "created_at" | "updated_at";
  sort_order?: "asc" | "desc";
  search?: string;
}

export async function getChatSessions(
  params: ChatHistoryParams = {}
): Promise<ApiResponse<ChatSessionsResponse>> {
  const searchParams = new URLSearchParams();
  const { sort_by = "updated_at", sort_order = "desc", search } = params;

  searchParams.append("sort_by", sort_by);
  searchParams.append("sort_order", sort_order);

  if (search) {
    searchParams.append("search", search);
  }

  const response = await axios.get<ApiResponse<ChatSessionsResponse>>(
    `${API_ENDPOINTS.chatHistory.base}${
      API_ENDPOINTS.chatHistory.sessions
    }?${searchParams.toString()}`
  );

  return response.data;
}

export async function getChatHistory(
  sessionId: string
): Promise<ApiResponse<ChatHistoryResponse>> {
  const response = await axios.get<ApiResponse<ChatHistoryResponse>>(
    `${API_ENDPOINTS.chatHistory.base}${API_ENDPOINTS.chatHistory.messages}/${sessionId}/messages`
  );

  return response.data;
}

export async function deleteChatSession(
  sessionId: string
): Promise<ApiResponse<{ success: boolean }>> {
  const response = await axios.delete<ApiResponse<{ success: boolean }>>(
    `${API_ENDPOINTS.chatHistory.base}${API_ENDPOINTS.chatHistory.sessions}/${sessionId}`
  );

  return response.data;
}
