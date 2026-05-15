"use client";
import { useEffect, useRef, useState } from "react";

export function VideoPlayer({ src, mediaId, episodeId, profileId, initialProgress = 0, nextHref }: { src: string; mediaId: string; episodeId?: string; profileId: string; initialProgress?: number; nextHref?: string }) {
  const ref = useRef<HTMLVideoElement>(null);
  const [status, setStatus] = useState("Bereit");
  useEffect(() => { if (ref.current && initialProgress > 0) ref.current.currentTime = initialProgress; }, [initialProgress]);
  useEffect(() => {
    const id = setInterval(async () => {
      const video = ref.current;
      if (!video || video.paused) return;
      await fetch("/api/watch/progress", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ mediaId, episodeId, profileId, progressSeconds: Math.floor(video.currentTime), completed: video.ended }) });
      setStatus(`Gespeichert bei ${Math.floor(video.currentTime)}s`);
    }, 15000);
    return () => clearInterval(id);
  }, [mediaId, episodeId, profileId]);
  return <div className="min-h-screen bg-black"><video ref={ref} src={src} controls playsInline className="h-[75vh] w-full bg-black" onEnded={() => nextHref && (window.location.href = nextHref)} /><div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4 text-sm text-zinc-300"><span>{status}</span><a href="/browse" className="rounded bg-white/10 px-4 py-2 hover:bg-white/20">Zur Übersicht</a></div></div>;
}
