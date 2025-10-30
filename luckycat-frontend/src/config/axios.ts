import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";
import { getAccessToken, getRefreshToken, clearAuthTokens } from "@/lib/cookie";
import { refreshToken } from "@/apis/auth";
import { API_TIMEOUT } from "./api";

interface QueueItem {
  resolve: (token: string) => void;
  reject: (error: unknown) => void;
}

interface ExtendedAxiosRequestConfig extends InternalAxiosRequestConfig {
  _retry?: boolean;
}

const axiosInstance = axios.create({
  timeout: API_TIMEOUT.default,
  headers: {
    "Content-Type": "application/json",
  },
});

let isRefreshing = false;
let failedQueue: QueueItem[] = [];

const processQueue = (
  error: unknown | null,
  token: string | null = null
): void => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else if (token) {
      resolve(token);
    }
  });

  failedQueue = [];
};

axiosInstance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const accessToken = getAccessToken();

    if (accessToken && config.headers) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }

    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as ExtendedAxiosRequestConfig;

    const shouldRefresh =
      (error.response?.status === 401 || error.response?.status === 403) &&
      originalRequest &&
      !originalRequest._retry;

    if (!shouldRefresh) {
      return Promise.reject(error);
    }

    if (isRefreshing) {
      return new Promise<string>((resolve, reject) => {
        failedQueue.push({ resolve, reject });
      })
        .then((token) => {
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`;
          }
          return axiosInstance(originalRequest);
        })
        .catch((err) => Promise.reject(err));
    }

    originalRequest._retry = true;
    isRefreshing = true;

    const refreshTokenValue = getRefreshToken();

    if (!refreshTokenValue) {
      clearAuthTokens();
      processQueue(new Error("No refresh token available"), null);
      isRefreshing = false;
      return Promise.reject(error);
    }

    try {
      await refreshToken({ refresh_token: refreshTokenValue });

      const newAccessToken = getAccessToken();

      if (!newAccessToken) {
        throw new Error("No access token received after refresh");
      }

      if (originalRequest.headers) {
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
      }

      processQueue(null, newAccessToken);
      isRefreshing = false;

      return axiosInstance(originalRequest);
    } catch (refreshError) {
      clearAuthTokens();
      processQueue(refreshError, null);
      isRefreshing = false;

      return Promise.reject(refreshError);
    }
  }
);

export default axiosInstance;
