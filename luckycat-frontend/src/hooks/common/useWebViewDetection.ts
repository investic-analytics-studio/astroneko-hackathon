import { useEffect, useState } from "react";
import { detectWebView, getCurrentUrl } from "@/lib/webviewDetection";

export function useWebViewDetection() {
  const [webViewInfo, setWebViewInfo] = useState({
    isWebView: false,
    userAgent: "",
    platformInfo: "",
  });
  const [currentUrl, setCurrentUrl] = useState<string>("");

  useEffect(() => {
    const info = detectWebView();
    const url = getCurrentUrl();

    setWebViewInfo(info);
    setCurrentUrl(url);
  }, []);

  return {
    isWebView: webViewInfo.isWebView,
    webViewInfo,
    currentUrl,
  };
}
