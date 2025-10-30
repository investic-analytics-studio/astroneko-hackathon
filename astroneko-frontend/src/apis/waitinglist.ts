import { API_ENDPOINTS } from "@/config/api";
import type { ApiResponse } from "@/types/api";
import axios from "../config/axios";

export interface WaitingListRequest {
  email: string;
}

export interface WaitingListData {
  id: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface WaitingListCheckData {
  success: boolean;
  message: string;
}

export async function addToWaitingList(
  params: WaitingListRequest
): Promise<ApiResponse<WaitingListData>> {
  const response = await axios.post<ApiResponse<WaitingListData>>(
    `${API_ENDPOINTS.waitlist.base}${API_ENDPOINTS.waitlist.join}`,
    params
  );

  return response.data;
}

export async function checkWaitingList(
  params: WaitingListRequest
): Promise<ApiResponse<WaitingListCheckData>> {
  const response = await axios.post<ApiResponse<WaitingListCheckData>>(
    `${API_ENDPOINTS.waitlist.base}${API_ENDPOINTS.waitlist.check}`,
    params
  );

  return response.data;
}
