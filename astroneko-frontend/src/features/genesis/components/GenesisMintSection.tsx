import { HelioCheckout } from "@heliofi/checkout-react";

const helioConfig = {
  paylinkId: "68ff3c94936eec0f079d5ed9", // prod : 68ff3c94936eec0f079d5ed9, dev : 68fb53162c788da645e0e591
  theme: {
    themeMode: "dark" as const,
  },
  platform: "magic_eden" as const,
  primaryColor: "#f7c36d",
  neutralColor: "#5A6578",
  backgroundColor: "#120c18",
  display: "inline" as const,
  onSuccess: (event: any) => console.log(event),
  onError: (event: any) => console.log(event),
  onPending: (event: any) => console.log(event),
  onCancel: () => console.log("Cancelled payment"),
  onStartPayment: () => console.log("Starting payment"),
};

export function GenesisMintSection() {
  return (
    <div className="flex flex-col">
      <div className="w-full rounded-[30px] overflow-hidden p-4 md:p-0">
        <div className="helio-checkout-wrapper w-full">
          <HelioCheckout config={helioConfig} />
        </div>
      </div>
        <div
          className="mt-4 w-[400px] mx-auto rounded-xl border border-[#f7c36d]/30 bg-[#f7c36d]/5 p-4 text-sm text-[#f7c36d]"
          role="note"
          aria-live="polite"
        >
          <p className="font-medium">Presale Info</p>
          <ul className="mt-2 list-disc space-y-1 pl-5 text-[#f7c36d]">
            <li>Presale participants will receive the NFT via airdrop after the presale ends.</li>
            <li>You will also receive a notification via email when your airdrop is sent.</li>
          </ul>
        </div>
    </div>
  );
}
