import { ReactNode } from "react";

interface ChatLayoutProps {
  chatContent: ReactNode;
  chatInput: ReactNode;
  chatRef?: React.RefObject<HTMLDivElement | null>;
}

export const ChatLayout = ({
  chatContent,
  chatInput,
  chatRef,
}: ChatLayoutProps) => (
  <div className="w-full md:h-screen flex flex-col bg-[image:var(--bg-display-2)] bg-cover bg-center text-white font-sans animate-fade-in overflow-hidden">
    <div
      ref={chatRef}
      className="mt-16 sm:mt-12 md:mt-0 w-full flex-1 overflow-y-scroll overflow-x-hidden px-0 sm:px-4 md:px-6 lg:px-8 xl:px-12 py-4 space-y-4 font-sans pt-20 sm:pt-24 md:pt-28 lg:pt-32 xl:pt-36 2xl:pt-40 pb-32 xl:min-w-[1100px] max-w-[1100px] mx-auto scroll-smooth transition-all duration-300 ease-in-out scrollbar-hide"
    >
      {chatContent}
    </div>
    <div className="w-full xl:min-w-[1100px] max-w-[1100px] mx-auto sticky bottom-0 transition-all duration-300 ease-in-out">
      {chatInput}
    </div>
  </div>
);
