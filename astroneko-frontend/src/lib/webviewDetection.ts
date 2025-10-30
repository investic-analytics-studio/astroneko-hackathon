import Bowser from "bowser";
import inAppSpy from "inapp-spy";
import { logger } from "./logger";

export interface WebViewInfo {
  isWebView: boolean;
  userAgent: string;
  platformInfo: string;
  isAndroid: boolean;
  isIOS: boolean;
  browser: string;
}

export const detectWebView = (): WebViewInfo => {
  const userAgent = navigator.userAgent;
  const browser = Bowser.getParser(userAgent);
  const platformInfo = browser.getOSName();

  // Use Bowser for reliable platform detection
  const isAndroid = browser.getOSName() === "Android";
  const isIOS = browser.getOSName() === "iOS";

  // Check if running in crypto wallet browsers - these are NOT webviews
  const isOKXWallet = /OKApp/i.test(userAgent);
  const isMetaMask = /MetaMask/i.test(userAgent);
  const isTrustWallet = /Trust/i.test(userAgent);
  const isCoinbaseWallet = /CoinbaseBrowser/i.test(userAgent);
  const isPhantomWallet = /Phantom/i.test(userAgent);
  const isBinanceWallet = /BinanceWallet/i.test(userAgent);
  const isRainbowWallet = /Rainbow/i.test(userAgent);

  const isCryptoWallet =
    isOKXWallet ||
    isMetaMask ||
    isTrustWallet ||
    isCoinbaseWallet ||
    isPhantomWallet ||
    isBinanceWallet ||
    isRainbowWallet;

  // Check if running in social media app browsers - these ARE webviews
  const isFacebookApp = /FBAN|FBAV|FB_IAB/i.test(userAgent);
  const isInstagramApp = /Instagram/i.test(userAgent);
  const isTwitterApp = /Twitter/i.test(userAgent);
  const isXApp = /X\.com/i.test(userAgent);
  const isThreadsApp = /Barcelona/i.test(userAgent);
  const isLineApp = /Line/i.test(userAgent);
  const isTikTokApp = /TikTok|musical_ly/i.test(userAgent);
  const isLemon8App = /Lemon8/i.test(userAgent);
  const isMessengerApp = /Messenger/i.test(userAgent);
  const isWhatsAppApp = /WhatsApp/i.test(userAgent);
  const isWeChatApp = /MicroMessenger/i.test(userAgent);
  const isTelegramApp = /Telegram/i.test(userAgent);
  const isSnapchatApp = /Snapchat/i.test(userAgent);
  const isRedditApp = /Reddit/i.test(userAgent);
  const isLinkedInApp = /LinkedInApp/i.test(userAgent);
  const isPinterestApp = /Pinterest/i.test(userAgent);
  const isDiscordApp = /Discord/i.test(userAgent);

  const isSocialMediaApp =
    isFacebookApp ||
    isInstagramApp ||
    isTwitterApp ||
    isXApp ||
    isThreadsApp ||
    isLineApp ||
    isTikTokApp ||
    isLemon8App ||
    isMessengerApp ||
    isWhatsAppApp ||
    isWeChatApp ||
    isTelegramApp ||
    isSnapchatApp ||
    isRedditApp ||
    isLinkedInApp ||
    isPinterestApp ||
    isDiscordApp;

  // Determine if this is a webview:
  // - Crypto wallets are always treated as NOT webviews (return false)
  // - Social media apps are always treated as webviews (return true)
  const isWebView = isCryptoWallet ? false : isSocialMediaApp;

  return {
    isWebView,
    userAgent,
    platformInfo,
    isAndroid,
    isIOS,
    browser: browser.getBrowserName(),
  };
};

export const getCurrentUrl = (): string => {
  return window.location.href;
};

export const isAndroid = (): boolean => {
  const browser = Bowser.getParser(navigator.userAgent);
  return browser.getOSName() === "Android";
};

export const isiOS = (): boolean => {
  const browser = Bowser.getParser(navigator.userAgent);
  return browser.getOSName() === "iOS";
};

export const isiOS17Plus = (): boolean => {
  if (!isiOS()) return false;

  const browser = Bowser.getParser(navigator.userAgent);
  const osVersion = browser.getOSVersion();
  return parseFloat(osVersion) >= 17;
};

export const isInAppBrowser = (): boolean => {
  return inAppSpy().isInApp;
};

export const shouldShowExternalBrowserButton = (): boolean => {
  return isInAppBrowser() && (isAndroid() || isiOS());
};

export const getExternalBrowserUrl = (
  currentUrl: string = window.location.href
): string => {
  const url = new URL(currentUrl);

  if (isAndroid()) {
    // Add query parameter to prevent infinite loops
    url.searchParams.set("external", "true");
    return `intent://${url.host}${url.pathname}${url.search}#Intent;scheme=https;end`;
  }

  if (isiOS17Plus()) {
    return `x-safari-https://${url.host}${url.pathname}${url.search}`;
  }

  if (isiOS()) {
    // Fallback for older iOS versions using shortcuts method
    const fallbackUrl = encodeURIComponent(currentUrl);
    return `shortcuts://x-callback-url/run-shortcut?name=${crypto.randomUUID()}&x-error=${fallbackUrl}`;
  }

  return currentUrl;
};

export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    // First try the modern Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      try {
        await navigator.clipboard.writeText(text);
        return true;
      } catch (err) {
        logger.warn("Clipboard API failed, trying fallback:", err);
      }
    }

    // Fallback for browsers that don't support Clipboard API or in insecure contexts
    const textArea = document.createElement("textarea");
    textArea.value = text;

    // Styling to ensure the element is accessible but not visible
    textArea.style.position = "fixed";
    textArea.style.top = "50%";
    textArea.style.left = "50%";
    textArea.style.transform = "translate(-50%, -50%)";
    textArea.style.width = "300px";
    textArea.style.height = "50px";
    textArea.style.zIndex = "9999";
    textArea.style.opacity = "1";
    textArea.style.backgroundColor = "white";
    textArea.style.color = "black";
    textArea.style.border = "2px solid #007AFF";
    textArea.style.borderRadius = "8px";
    textArea.style.padding = "8px";
    textArea.style.fontSize = "16px"; // Prevents zoom on iOS
    textArea.setAttribute("readonly", "readonly");

    document.body.appendChild(textArea);

    // Focus and select the text
    textArea.focus();
    textArea.select();
    textArea.setSelectionRange(0, text.length);

    // Add a small delay for better compatibility
    await new Promise((resolve) => setTimeout(resolve, 100));

    let success = false;
    try {
      // Try to copy using the deprecated execCommand as last resort
      // This will show a deprecation warning but still works in most browsers
      success = document.execCommand("copy");
    } catch (err) {
      logger.warn("Legacy copy method failed:", err);
      // If execCommand fails, we can try prompting the user
      success = false;
    }

    // Clean up after a short delay to ensure copy operation completes
    setTimeout(() => {
      if (document.body.contains(textArea)) {
        document.body.removeChild(textArea);
      }
    }, 500);

    // If all methods fail, we could show a manual copy prompt
    if (!success) {
      logger.warn(
        "All copy methods failed. Consider showing manual copy instructions to user."
      );
    }

    return success;
  } catch (error) {
    logger.error("Copy to clipboard failed:", error);
    return false;
  }
};
