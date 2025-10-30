import { ChevronDown, ChevronUp, PawPrint } from "lucide-react";

interface ChatSuggestionsProps {
  showQuestions: boolean;
  onToggleQuestions: () => void;
  onSelectQuestion: (question: string) => void;
  questions: string[];
}

export const ChatSuggestions = ({
  showQuestions,
  onToggleQuestions,
  onSelectQuestion,
  questions,
}: ChatSuggestionsProps) => {
  return (
    <div className="relative mb-2">
      {showQuestions && (
        <div className="absolute bottom-[30px] 2xl:bottom-[110px] left-0 w-full px-2 sm:px-3 md:px-4 py-2 pt-4 sm:pt-5 md:pt-6 overflow-hidden animate-slide-down">
          <div className="max-w-4xl 2xl:max-w-[1700px] mx-auto grid grid-cols-1 sm:grid-cols-2 2xl:grid-cols-2 gap-2 2xl:gap-5">
            {Array.isArray(questions)
              ? questions.map((question) => (
                  <button
                    key={question}
                    onClick={() => onSelectQuestion(question)}
                    className="group relative px-2 sm:px-3 md:px-4 lg:px-4 2xl:px-5 text-left
                           bg-[#4F3533] font-normal
                           hover:from-white/20 hover:to-white/20
                           border border-white/20 hover:border-white/30
                           rounded-xl 2xl:rounded-[24px] text-white
                           transition-all duration-300"
                  >
                    <div className="flex items-center gap-3">
                      {/* <Sparkles className="w-4 2xl:w-6 h-4 2xl:h-6 text-[#209CFF] group-hover:opacity-100 transition-opacity" /> */}
                      <div className="flex-shrink-0">
                        <PawPrint className="w-6 h-6 md:w-7 md:h-7 lg:w-8 lg:h-8 2xl:w-10 2xl:h-10 text-[#F7C36D] fill-current" />
                      </div>
                      <span className="text-xs sm:text-sm md:text-sm lg:text-base 2xl:text-[30px] leading-tight overflow-hidden text-ellipsis line-clamp-2">
                        {question}
                      </span>
                    </div>
                    {/* Animated star trail */}
                    <div className="absolute inset-0 overflow-hidden rounded-lg animate-pulse-glow" />
                  </button>
                ))
              : null}
          </div>
        </div>
      )}
      <button
        onClick={onToggleQuestions}
        className="absolute w-auto h-[26px] md:h-[32px] lg:h-[40px] 2xl:w-[340px] 2xl:h-[80px] -top-6 md:-top-7 lg:-top-8 2xl:-top-[90px] left-1/2 -translate-x-1/2
                 px-3 sm:px-4 md:px-5 lg:px-5 2xl:px-6 py-1.5 flex items-center gap-1.5
                 bg-[#FFFFFF]/25 justify-center
                 hover:from-purple-500/30 hover:to-indigo-500/30
                 border border-white/10 hover:border-white/30
                 rounded-full text-white focus:outline-none focus:ring-none
                 text-[10px] sm:text-[11px] md:text-[12px] lg:text-[14px] 2xl:text-[24px] font-normal transition-all duration-300"
      >
        {showQuestions ? (
          <>
            Hide suggestions{" "}
            <ChevronDown className="w-3 h-3 md:w-3.5 md:h-3.5 lg:w-4 lg:h-4 2xl:w-6 2xl:h-6" />
          </>
        ) : (
          <>
            Show suggestions{" "}
            <ChevronUp className="w-3 h-3 md:w-3.5 md:h-3.5 lg:w-4 lg:h-4 2xl:w-6 2xl:h-6" />
          </>
        )}
      </button>
    </div>
  );
};
