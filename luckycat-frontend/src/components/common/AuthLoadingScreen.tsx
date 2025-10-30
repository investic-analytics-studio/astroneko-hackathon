export function AuthLoadingScreen() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-[image:var(--bg-display-2)] bg-cover bg-center">
      <div className="flex flex-col items-center gap-4">
        <div className="w-8 h-8 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
        <p className="text-white/70 text-sm">Loading...</p>
      </div>
    </div>
  );
}
