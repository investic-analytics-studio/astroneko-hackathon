import {
  GenesisDetailSection,
  GenesisImageSection,
  GenesisMintSection,
} from "@/features/genesis";

export default function GenesisPage() {
  return (
    <div className="bg-[#221A29] !h-auto items-center justify-center py-28 px-4 md:px-8">
      <div className="w-full max-w-7xl mx-auto opacity-0 animate-fade-in">
        <div className="mt-4 md:mt-0 flex flex-col-reverse lg:flex-row gap-4 lg:gap-0 items-start">
          {/* Left Column: Single Card with Image & Details */}
          <div className="w-full">
            <GenesisImageSection />
            <GenesisDetailSection />
          </div>

          {/* Right Column: Mint Section */}
          <div className="w-full flex justify-center">
            <GenesisMintSection />
          </div>
        </div>
      </div>
    </div>
  );
}
