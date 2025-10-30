import { API_ENDPOINTS } from "@/config/api";
import type { ApiResponse } from "@/types/api";
import axios from "../config/axios";

export interface UserLimitResponse {
  is_limit_reached: boolean;
  current_usage?: number;
  max_limit?: number;
  reset_date?: string;
}

export async function isUserLimitReached(): Promise<
  ApiResponse<UserLimitResponse>
> {
  const response = await axios.get<ApiResponse<UserLimitResponse>>(
    `${API_ENDPOINTS.userLimit.base}${API_ENDPOINTS.userLimit.check}`
  );

  return response.data;
}
