import React, { useState, useRef, useEffect } from "react";

interface AudioPlayerProps {
	src: string;
	title: string;
	artist: string;
	artworkUrl?: string;
	compact?: boolean;
}

export default function AudioPlayer({
	src,
	title,
	artist,
	artworkUrl,
	compact = false,
}: AudioPlayerProps) {
	const audioRef = useRef<HTMLAudioElement>(null);
	const [isPlaying, setIsPlaying] = useState(false);
	const [currentTime, setCurrentTime] = useState(0);
	const [duration, setDuration] = useState(0);

	useEffect(() => {
		const audio = audioRef.current;
		if (!audio) return;

		const onTimeUpdate = () => setCurrentTime(audio.currentTime);
		const onLoadedMetadata = () => setDuration(audio.duration);
		const onEnded = () => setIsPlaying(false);

		audio.addEventListener("timeupdate", onTimeUpdate);
		audio.addEventListener("loadedmetadata", onLoadedMetadata);
		audio.addEventListener("ended", onEnded);

		return () => {
			audio.removeEventListener("timeupdate", onTimeUpdate);
			audio.removeEventListener("loadedmetadata", onLoadedMetadata);
			audio.removeEventListener("ended", onEnded);
		};
	}, []);

	const togglePlay = () => {
		const audio = audioRef.current;
		if (!audio) return;

		if (isPlaying) {
			audio.pause();
		} else {
			audio.play();
		}
		setIsPlaying(!isPlaying);
	};

	const seek = (e: React.MouseEvent<HTMLDivElement>) => {
		const audio = audioRef.current;
		if (!audio || !duration) return;

		const rect = e.currentTarget.getBoundingClientRect();
		const percent = (e.clientX - rect.left) / rect.width;
		audio.currentTime = percent * duration;
	};

	const formatTime = (t: number) => {
		const mins = Math.floor(t / 60);
		const secs = Math.floor(t % 60);
		return `${mins}:${secs.toString().padStart(2, "0")}`;
	};

	const progress = duration ? (currentTime / duration) * 100 : 0;

	if (compact) {
		return (
			<div className="flex items-center gap-3">
				<audio ref={audioRef} src={src} preload="metadata" />
				<button
					onClick={togglePlay}
					className="w-10 h-10 rounded-full bg-muziboo-gold text-muziboo-bg flex items-center justify-center hover:bg-muziboo-text transition-colors flex-shrink-0"
					aria-label={isPlaying ? "Pause" : "Play"}
				>
					{isPlaying ? (
						<svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
							<rect x="6" y="4" width="4" height="16" />
							<rect x="14" y="4" width="4" height="16" />
						</svg>
					) : (
						<svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
							<polygon points="5,3 19,12 5,21" />
						</svg>
					)}
				</button>
				<div className="flex-1 min-w-0">
					<div
						className="h-1.5 bg-muziboo-bg rounded-full cursor-pointer group"
						onClick={seek}
					>
						<div
							className="h-full bg-muziboo-gold rounded-full transition-all duration-100 group-hover:bg-muziboo-text"
							style={{ width: `${progress}%` }}
						/>
					</div>
				</div>
				<span className="text-xs text-muziboo-text-muted tabular-nums flex-shrink-0">
					{formatTime(currentTime)}
				</span>
			</div>
		);
	}

	return (
		<div className="bg-muziboo-bg-softer rounded-xl p-4 border border-muziboo-border">
			<audio ref={audioRef} src={src} preload="metadata" />

			<div className="flex items-center gap-4">
				{artworkUrl && (
					<img
						src={artworkUrl}
						alt={`${title} artwork`}
						className="w-14 h-14 rounded-lg object-cover flex-shrink-0"
					/>
				)}

				<div className="flex-1 min-w-0">
					<p className="text-muziboo-text font-medium truncate">{title}</p>
					<p className="text-muziboo-text-muted text-sm truncate">{artist}</p>
				</div>

				<button
					onClick={togglePlay}
					className="w-12 h-12 rounded-full bg-muziboo-gold text-muziboo-bg flex items-center justify-center hover:bg-muziboo-text transition-colors flex-shrink-0"
					aria-label={isPlaying ? "Pause" : "Play"}
				>
					{isPlaying ? (
						<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
							<rect x="6" y="4" width="4" height="16" />
							<rect x="14" y="4" width="4" height="16" />
						</svg>
					) : (
						<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
							<polygon points="5,3 19,12 5,21" />
						</svg>
					)}
				</button>
			</div>

			<div className="mt-3 flex items-center gap-3">
				<span className="text-xs text-muziboo-text-muted tabular-nums w-10 text-right">
					{formatTime(currentTime)}
				</span>
				<div
					className="flex-1 h-1.5 bg-muziboo-bg rounded-full cursor-pointer group"
					onClick={seek}
				>
					<div
						className="h-full bg-muziboo-gold rounded-full transition-all duration-100 group-hover:bg-muziboo-text"
						style={{ width: `${progress}%` }}
					/>
				</div>
				<span className="text-xs text-muziboo-text-muted tabular-nums w-10">
					{formatTime(duration)}
				</span>
			</div>
		</div>
	);
}
