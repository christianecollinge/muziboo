import React from "react";
import TrackCard from "./TrackCard";

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

interface TrackGridProps {
	tracks: Track[];
	emptyMessage?: string;
}

export default function TrackGrid({
	tracks,
	emptyMessage = "No tracks yet",
}: TrackGridProps) {
	if (tracks.length === 0) {
		return (
			<div className="text-center py-20">
				<svg
					width="48"
					height="48"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					strokeWidth="1.5"
					className="mx-auto text-muziboo-text-muted/30 mb-4"
				>
					<path d="M9 18V5l12-2v13" />
					<circle cx="6" cy="18" r="3" />
					<circle cx="18" cy="16" r="3" />
				</svg>
				<p className="text-muziboo-text-muted">{emptyMessage}</p>
			</div>
		);
	}

	return (
		<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
			{tracks.map((track) => (
				<TrackCard key={track.id} track={track} />
			))}
		</div>
	);
}
