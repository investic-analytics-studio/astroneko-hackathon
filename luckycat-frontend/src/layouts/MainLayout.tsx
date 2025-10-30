import { useLocation } from "@tanstack/react-router";
import { Header } from "@/components/common/Header";
import { ChatHistorySidebar } from "@/features/chat/history";
import { useAppNavigation, useWebViewDetection } from "@/hooks";
import { WebViewWarningDialog } from "@/components/common/WebViewWarningDialog";

interface MainLayoutProps {
  children: React.ReactNode;
  showHeader?: boolean;
}

function PageWrapper({ children }: { children: React.ReactNode }) {
  return (
    <div
      className="page-wrapper opacity-0 animate-fade-in"
      style={{ backgroundColor: "rgba(69, 10, 10, 0.6)" }}
    >
      {children}
    </div>
  );
}

export function MainLayout({ children, showHeader = true }: MainLayoutProps) {
  const { goToHome } = useAppNavigation();
  const { isWebView, currentUrl } = useWebViewDetection();
  const location = useLocation();

  // Check if current page is a chat page
  const isChatPage = location.pathname.includes("/chat/");

  return (
    <div className="w-full min-h-screen bg-[image:var(--bg-display-2)] bg-cover bg-center text-white overflow-hidden">
      {/* Header and Children */}
      <div className="flex flex-col h-screen">
        {showHeader && <Header onScrollToHome={goToHome} />}

        <div className="flex-1 overflow-auto">
          <PageWrapper>{children}</PageWrapper>
        </div>
      </div>

      {/* Chat History Sidebar - only show on chat pages */}
      {isChatPage && <ChatHistorySidebar />}

      {/* WebView Warning Dialog - appears globally when in webview */}
      <WebViewWarningDialog isOpen={isWebView} currentUrl={currentUrl} />
    </div>
  );
}
