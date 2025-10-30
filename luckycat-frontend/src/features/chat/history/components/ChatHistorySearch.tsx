import { Search, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { memo } from "react";

interface ChatHistorySearchProps {
  value: string;
  onChange: (value: string) => void;
}

export const ChatHistorySearch = memo(
  ({ value, onChange }: ChatHistorySearchProps) => {
    const handleClear = () => {
      onChange("");
    };

    return (
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-white/50 pointer-events-none" />
        <Input
          type="text"
          placeholder="Search conversations..."
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="pl-10 pr-10 bg-white/10 border-white/20 text-white placeholder:text-white/50 focus:border-[var(--brand-accent)] focus:ring-[var(--brand-accent)]/20"
        />
        {value && (
          <Button
            variant="ghost"
            size="icon"
            onClick={handleClear}
            className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 text-white/50 hover:text-white hover:bg-white/10"
          >
            <X className="h-3 w-3" />
            <span className="sr-only">Clear search</span>
          </Button>
        )}
      </div>
    );
  }
);

ChatHistorySearch.displayName = "ChatHistorySearch";
