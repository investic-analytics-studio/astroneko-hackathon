import Cookies from "js-cookie";

const SESSION_ID_KEY = "session_id";
const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";

export const setSessionId = (sessionId: string, category: string = "general"): void => {
  const sessionKey = `${SESSION_ID_KEY}_${category}`;
  Cookies.set(sessionKey, sessionId, {
    expires: 7, // Cookie expires in 7 days
    sameSite: "strict",
  });
};

export const getSessionId = (category: string = "general"): string | undefined => {
  const sessionKey = `${SESSION_ID_KEY}_${category}`;
  return Cookies.get(sessionKey);
};

export const removeSessionId = (category: string = "general"): void => {
  const sessionKey = `${SESSION_ID_KEY}_${category}`;
  Cookies.remove(sessionKey);
};

export const removeAllSessionIds = (): void => {
  // Remove session IDs for all categories
  const categories = ["general", "crypto", "lover"];
  categories.forEach(category => {
    const sessionKey = `${SESSION_ID_KEY}_${category}`;
    Cookies.remove(sessionKey);
  });
  // Also remove the old session_id key for backward compatibility
  Cookies.remove(SESSION_ID_KEY);
};

export const setAccessToken = (token: string, expiresIn: number): void => {
  const expirationDate = new Date();
  expirationDate.setSeconds(expirationDate.getSeconds() + expiresIn);

  Cookies.set(ACCESS_TOKEN_KEY, token, {
    expires: expirationDate,
    sameSite: "strict",
    secure: true,
  });
};

export const getAccessToken = (): string | undefined => {
  return Cookies.get(ACCESS_TOKEN_KEY);
};

export const setRefreshToken = (token: string): void => {
  Cookies.set(REFRESH_TOKEN_KEY, token, {
    expires: 30, // Refresh token expires in 30 days
    sameSite: "strict",
    secure: true,
  });
};

export const getRefreshToken = (): string | undefined => {
  return Cookies.get(REFRESH_TOKEN_KEY);
};

export const clearAuthTokens = (): void => {
  Cookies.remove(ACCESS_TOKEN_KEY);
  Cookies.remove(REFRESH_TOKEN_KEY);
  removeAllSessionIds();
};
