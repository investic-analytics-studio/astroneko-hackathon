import { ArrowLeft } from "lucide-react";

interface BackButtonProps {
  onBack: () => void;
}

export function BackButton({ onBack }: BackButtonProps) {
  return (
    <div className="flex justify-start mb-8 max-w-7xl mx-auto w-full opacity-0 -translate-x-5 animate-fade-in">
      <button
        onClick={onBack}
        className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 hover:border-none focus:border-none focus:ring-0 focus:outline-none text-white rounded-full backdrop-blur-sm transition-all duration-300 group hover:scale-105 active:scale-95"
      >
        <ArrowLeft className="w-5 h-5 group-hover:-translate-x-1 transition-transform duration-300" />
        <span className="font-medium">Back</span>
      </button>
    </div>
  );
}
