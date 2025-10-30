export function GenesisDetailSection() {
  return (
    <div className="w-full mt-8 pt-8 border-t border-[#2A2A2A] bg-[#120C18] rounded-[30px] p-8">
      {/* Collection Title */}
      <h1 className="text-[28px] font-bold text-white mb-3 text-center">
        AstroNeko Genesis Collection
      </h1>

      {/* Description */}
      <p className="text-[14px] text-[#A1A1AA] leading-relaxed mb-4">
        The Genesis Collection is a AstroNeko masterfully consecrated set,
        blessed under the stars to infuse its owner with cosmic energy for
        wealth and wisdom.
      </p>

      {/* limited card  */}
      <div className="mb-4">
        <div className="text-[24px] font-bold text-[#f7c36d]">
          50 Total - 33 Presale
        </div>
        <div className="text-[14px] text-[#f7c36d]">at $2,989 USD</div>
      </div>

      {/* Collection Info */}
      <p className="text-[14px] text-[#A1A1AA] leading-relaxed border-b border-[#2A2A2A] pb-2 mb-2">
        This exclusive tier is reserved for Founding Members, Angel Investors,
        Early Supporters, and True Believers who dare to hold a piece of destiny
        itself.
      </p>

      <div className="flex flex-col md:flex-row gap-6">
        {/* What's Inside Section */}
        <div className="flex-1 mb-6 md:mb-0">
          <h3 className="text-[18px] font-bold text-white mb-3">
            What's Inside
          </h3>
          <div className="space-y-2">
            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                AstroNeko Genesis NFT
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Onchain Solana NFT , Digital passport for Astro Star Inner
                Circle
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                {`Physical Collector’s Boxset`}
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Sacred Precision-crafted case for your exclusive relics
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                {`Major Arcana Tarot (22 Cards)`}
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                {`Genesis Tarot 22-card deck by AstroNeko’s supreme star`}
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                Crypto Orb & Bracelet{" "}
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Astroneko Oracle Orb & Bracelet talisman for wealth
              </div>
            </div>
          </div>
        </div>

        {/* Genesis Benefits */}
        <div className="flex-1 md:pt-0 md:border-t-0 md:border-l md:pl-6 md:ml-6 border-t border-[#2A2A2A] pt-6">
          <h3 className="text-[18px] font-bold text-white mb-3">
            Genesis Benefits
          </h3>
          <div className="space-y-2">
            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                Future Airdrop Priority
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Exclusive access to upcoming airdrops. Stake your Genesis NFT to
                secure top-tier rewards
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                Premium AI Portal Access
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Unlock advanced tools across AstroTrading, AstroFight, and
                AstroFengshui — where data meets destiny
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                🔮 1-on-1 Astrologer Session
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                {`Complimentary consultation with AstroNeko’s star seer`}
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                Astro Star Inner Circle
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Join an inner circle network of members. Gain entry to private
                events and early product access
              </div>
            </div>

            <div className="bg-[#1A1523] rounded-lg p-3">
              <div className="text-[15px] font-semibold text-white mb-1">
                Astro Academy Free Enrollment
              </div>
              <div className="text-[12px] text-[#A1A1AA]">
                Early full access to Astro curated courses in trading , AI, and
                Astro intelligence{" "}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
