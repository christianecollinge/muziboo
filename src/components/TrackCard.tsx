import React, { useState, useRef } from "react";
import AudioPlayer from "./AudioPlayer";

interface Track {
	id: string;
	title: string;
	description?: string;
	genre?: string;
	audio_url: string;
	artwork_url?: string;
	created_at: string;
	profiles?: {
		username: string;
		display_name: string;
		avatar_url?: string;
	};
}

interface TrackCardProps {
	track: Track;
}

export default function TrackCard({ track }: TrackCardProps) {
	const [showPlayer, setShowPlayer] = useState(false);
	const profile = track.profiles;

	const timeAgo = (date: string) => {
		const seconds = Math.floor((Date.now() - new Date(date).getTime()) / 1000);
		if (seconds < 60) return "just now";
		const mins = Math.floor(seconds / 60);
		if (mins < 60) return `${mins}m ago`;
		const hours = Math.floor(mins / 60);
		if (hours < 24) return `${hours}h ago`;
		const days = Math.floor(hours / 24);
		if (days < 30) return `${days}d ago`;
		const months = Math.floor(days / 30);
		return `${months}mo ago`;
	};

	return (
		<div className="group rounded-xl overflow-hidden bg-muziboo-bg-softer border border-muziboo-border hover:border-muziboo-gold/40 transition-all duration-300">
			{/* Artwork */}
			<div
				className="relative aspect-square overflow-hidden cursor-pointer"
				onClick={() => setShowPlayer(!showPlayer)}
			>
				{track.artwork_url ? (
					<img
						src={track.artwork_url}
						alt={`${track.title} artwork`}
						className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
						loading="lazy"
					/>
				) : (
					<div className="w-full h-full bg-gradient-to-br from-muziboo-bar-1/30 via-muziboo-bar-3/30 to-muziboo-bar-5/30 flex items-center justify-center">
						<svg
							width="48"
							height="48"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							strokeWidth="1.5"
							className="text-muziboo-text-muted/40"
						>
							<path d="M9 18V5l12-2v13" />
							<circle cx="6" cy="18" r="3" />
							<circle cx="18" cy="16" r="3" />
						</svg>
					</div>
				)}

				{/* Play overlay */}
				<div className="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-colors duration-300 flex items-center justify-center">
					<div className="w-14 h-14 rounded-full bg-muziboo-gold text-muziboo-bg flex items-center justify-center opacity-0 group-hover:opacity-100 transform scale-90 group-hover:scale-100 transition-all duration-300 shadow-lg">
						<svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
							<polygon points="5,3 19,12 5,21" />
						</svg>
					</div>
				</div>

				{/* Genre badge */}
				{track.genre && (
					<span className="absolute top-3 left-3 px-2.5 py-1 text-xs font-medium rounded-full bg-black/60 text-muziboo-text backdrop-blur-sm">
						{track.genre}
					</span>
				)}
			</div>

			{/* Info */}
			<div className="p-4">
				<h3 className="font-semibold text-muziboo-text truncate mb-1">
					{track.title}
				</h3>

				{profile && (
					<a
						href={`/user/${profile.username}`}
						className="flex items-center gap-2 group/artist"
					>
						{profile.avatar_url ? (
							<img
								src={profile.avatar_url}
								alt={profile.display_name}
								className="w-5 h-5 rounded-full object-cover"
							/>
						) : (
							<div className="w-5 h-5 rounded-full bg-muziboo-gold/20 flex items-center justify-center">
								<span className="text-xs text-muziboo-gold font-medium">
									{(profile.display_name || profile.username)
										.charAt(0)
										.toUpperCase()}
								</span>
							</div>
						)}
						<span className="text-sm text-muziboo-text-muted group-hover/artist:text-muziboo-gold transition-colors truncate">
							{profile.display_name || profile.username}
						</span>
					</a>
				)}

				<div className="flex items-center justify-between mt-3">
					<span className="text-xs text-muziboo-text-muted/60">
						{timeAgo(track.created_at)}
					</span>
				</div>

				{/* Inline player */}
				{showPlayer && (
					<div className="mt-3 pt-3 border-t border-muziboo-border">
						<AudioPlayer
							src={track.audio_url}
							title={track.title}
							artist={profile?.display_name || "Unknown"}
							compact
						/>
					</div>
				)}
			</div>
		</div>
	);
}
