import { ChatInput } from "../core/ChatInput";
import { ChatSuggestions } from "../core/ChatSuggestions";

interface ChatInputAreaProps {
  input: string;
  loading: boolean;
  showSuggestions: boolean;
  questions: string[];
  onInputChange: (value: string) => void;
  onSendMessage: () => void;
  onToggleQuestions: () => void;
  onSelectQuestion: (question: string) => void;
}

export const ChatInputArea = ({
  input,
  loading,
  showSuggestions,
  questions,
  onInputChange,
  onSendMessage,
  onToggleQuestions,
  onSelectQuestion,
}: ChatInputAreaProps) => (
  <div className="flex flex-col gap-1 px-2 sm:px-4 md:px-6 lg:px-8 xl:px-12 2xl:px-16 pb-16 sm:pb-20 md:pb-24 lg:pb-28 xl:pb-32 2xl:pb-36 transition-all duration-300 ease-in-out animate-slide-in">
    <div className="transition-all duration-300 ease-in-out transform hover:scale-[1.01]">
      <ChatSuggestions
        showQuestions={showSuggestions}
        onToggleQuestions={onToggleQuestions}
        onSelectQuestion={onSelectQuestion}
        questions={questions}
      />
    </div>
    <div className="transition-all duration-300 ease-in-out transform hover:scale-[1.02] animate-fade-in">
      <ChatInput
        input={input}
        loading={loading}
        onInputChange={onInputChange}
        onSendMessage={onSendMessage}
      />
    </div>
  </div>
);