import { ChatMessage } from "./ChatMessage";
import { FortuneMessage } from "@/hooks/fortune/useFortune";
import { Loader2 } from "lucide-react";

interface ChatMessagesProps {
  messages: FortuneMessage[];
  isLoading?: boolean;
}

export const ChatMessages = ({ messages, isLoading }: ChatMessagesProps) => (
  <>
    {isLoading && (
      <div className="flex-1 flex items-center justify-center px-4 sm:px-6 md:px-8 lg:px-12 xl:px-16 animate-fade-in transition-all duration-300 ease-in-out">
        <div className="flex items-center gap-3">
          <Loader2 className="h-6 w-6 text-white/60 animate-spin" />
          <span className="text-white/60 font-press-start text-sm">
            Loading conversation...
          </span>
        </div>
      </div>
    )}
    <div className="space-y-4 transition-all duration-300 ease-in-out">
      {messages.map((message) => (
        <div key={message.id}>
          <ChatMessage message={message} />
        </div>
      ))}
    </div>
  </>
);
