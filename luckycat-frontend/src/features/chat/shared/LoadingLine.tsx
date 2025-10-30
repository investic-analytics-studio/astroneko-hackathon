export const LoadingTextLine = ({ width }: { width: string }) => (
  <div
    className="relative h-4 bg-white/20 rounded overflow-hidden animate-pulse"
    style={{ width }}
  >
    <div className="absolute inset-0 bg-gradient-to-r from-white/20 via-white/20 to-white/20 animate-shimmer" />
  </div>
);
