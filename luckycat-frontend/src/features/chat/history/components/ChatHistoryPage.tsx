import { ChatHistorySidebar } from "./ChatHistorySidebar";
import { ChatHistoryViewer } from "./ChatHistoryViewer";
import {
  useChatHistoryStore,
  chatHistorySelectors,
} from "@/features/chat/history";

export function ChatHistoryPage() {
  const currentSessionId = useChatHistoryStore(
    chatHistorySelectors.currentSessionId
  );

  return (
    <div className="w-full min-h-screen flex bg-[image:var(--bg-display-2)] bg-cover bg-center text-white font-sans">
      {/* Chat History Sidebar */}
      <ChatHistorySidebar />

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col">
        {currentSessionId ? (
          <ChatHistoryViewer className="flex-1" />
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <h2 className="text-white/60 font-press-start text-xl mb-4">
                Welcome to Chat History
              </h2>
              <p className="text-white/40 font-press-start text-sm">
                Select a conversation from the sidebar to view your chat history
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
