import { Button } from "@/components/ui/button";

interface FlowSelectionProps {
  onSelectReferral: () => void;
  onSelectWaitlist: () => void;
}

export function FlowSelection({
  onSelectReferral,
  onSelectWaitlist,
}: FlowSelectionProps) {
  return (
    <div className="space-y-6 mt-2">
      <div className="space-y-3">
        <Button
          onClick={onSelectReferral}
          className="w-full bg-[#F7C36D] text-black font-semibold border-none rounded-[12px]
                   py-4 h-auto hover:bg-[#FFB53A] hover:border-none focus:border-none focus:ring-0 focus:outline-none transition-all duration-300"
        >
          <div className="flex flex-col items-center gap-1">
            <span className="text-[16px] font-bold">
              I have a Referral code
            </span>
            <span className="text-[14px] font-medium">
              Get immediate access
            </span>
          </div>
        </Button>

        <Button
          onClick={onSelectWaitlist}
          variant="outline"
          className="w-full bg-white/10 py-4 h-auto border-white/10 rounded-[12px]
          hover:bg-white/20 hover:border-white/20 focus:border-none focus:ring-0 focus:outline-none transition-all duration-300"
        >
          <div className="flex flex-col items-center gap-1 text-white">
            <span className="text-[16px] font-bold">Join the waitlist</span>
            <span className="text-[14px] font-medium opacity-70">
              We'll email you when ready
            </span>
          </div>
        </Button>
      </div>
    </div>
  );
}
