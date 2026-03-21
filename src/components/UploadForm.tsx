import React, { useState, useRef } from "react";
import { getAccessToken } from "../lib/supabase";
import { uploadFile, api } from "../lib/api";

export default function UploadForm() {
	const [title, setTitle] = useState("");
	const [description, setDescription] = useState("");
	const [genre, setGenre] = useState("");
	const [audioFile, setAudioFile] = useState<File | null>(null);
	const [artworkFile, setArtworkFile] = useState<File | null>(null);
	const [artworkPreview, setArtworkPreview] = useState("");
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState("");
	const [success, setSuccess] = useState(false);
	const [uploadProgress, setUploadProgress] = useState("");

	const audioInputRef = useRef<HTMLInputElement>(null);
	const artworkInputRef = useRef<HTMLInputElement>(null);

	const handleArtworkChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const file = e.target.files?.[0];
		if (file) {
			setArtworkFile(file);
			const reader = new FileReader();
			reader.onload = (ev) => setArtworkPreview(ev.target?.result as string);
			reader.readAsDataURL(file);
		}
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError("");
		setLoading(true);
		setSuccess(false);

		try {
			const token = await getAccessToken();
			if (!token) {
				window.location.href = "/login";
				return;
			}

			if (!audioFile) {
				setError("Please select an audio file");
				setLoading(false);
				return;
			}

			// 1. Upload audio
			setUploadProgress("Uploading audio...");
			const audioResult = await uploadFile(
				"/api/upload/audio",
				audioFile,
				token
			);

			// 2. Upload artwork (optional)
			let artworkUrl = "";
			if (artworkFile) {
				setUploadProgress("Uploading artwork...");
				const artworkResult = await uploadFile(
					"/api/upload/artwork",
					artworkFile,
					token
				);
				artworkUrl = artworkResult.url;
			}

			// 3. Create track record
			setUploadProgress("Creating track...");
			await api("/api/tracks", {
				method: "POST",
				token,
				body: {
					title,
					description,
					genre,
					audio_url: audioResult.url,
					artwork_url: artworkUrl,
				},
			});

			setSuccess(true);
			setTitle("");
			setDescription("");
			setGenre("");
			setAudioFile(null);
			setArtworkFile(null);
			setArtworkPreview("");
			setUploadProgress("");
		} catch (err: any) {
			setError(err.message || "Upload failed");
			setUploadProgress("");
		} finally {
			setLoading(false);
		}
	};

	const genres = [
		"Rock",
		"Jazz",
		"Classical",
		"Electronic",
		"Hip Hop",
		"Folk",
		"Blues",
		"Metal",
		"Punk",
		"Ambient",
		"Experimental",
		"World",
		"Singer-Songwriter",
		"Lo-Fi",
		"Other",
	];

	return (
		<div className="w-full max-w-2xl mx-auto">
			<div className="bg-muziboo-bg-softer rounded-xl p-8 border border-muziboo-border">
				<h1 className="text-3xl font-semibold text-muziboo-text mb-2">
					Upload a Track
				</h1>
				<p className="text-muziboo-text-muted mb-8">
					Share what you're working on. Demos, fragments, works in progress —
					everything is welcome.
				</p>

				{success && (
					<div className="mb-6 bg-muziboo-teal/10 border border-muziboo-teal/30 rounded-lg px-4 py-3">
						<p className="text-muziboo-teal-light font-medium">
							Track uploaded!
						</p>
						<p className="text-muziboo-text-muted text-sm mt-1">
							Your track is now live.{" "}
							<a
								href="/dashboard"
								className="text-muziboo-gold hover:underline"
							>
								View your dashboard
							</a>
							{" · "}
							<button
								onClick={() => setSuccess(false)}
								className="text-muziboo-gold hover:underline"
							>
								Upload another
							</button>
						</p>
					</div>
				)}

				<form onSubmit={handleSubmit} className="space-y-6">
					{/* Audio file */}
					<div>
						<label className="block text-sm font-medium text-muziboo-text mb-2">
							Audio File *
						</label>
						<input
							ref={audioInputRef}
							type="file"
							accept="audio/*"
							onChange={(e) => setAudioFile(e.target.files?.[0] || null)}
							className="hidden"
						/>
						<button
							type="button"
							onClick={() => audioInputRef.current?.click()}
							className="w-full border-2 border-dashed border-muziboo-border rounded-xl py-8 px-4 text-center hover:border-muziboo-gold/50 transition-colors group"
						>
							{audioFile ? (
								<div>
									<svg
										width="32"
										height="32"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										strokeWidth="1.5"
										className="mx-auto text-muziboo-gold mb-2"
									>
										<path d="M9 18V5l12-2v13" />
										<circle cx="6" cy="18" r="3" />
										<circle cx="18" cy="16" r="3" />
									</svg>
									<p className="text-muziboo-text font-medium">
										{audioFile.name}
									</p>
									<p className="text-muziboo-text-muted text-sm mt-1">
										{(audioFile.size / (1024 * 1024)).toFixed(1)} MB · Click to
										change
									</p>
								</div>
							) : (
								<div>
									<svg
										width="32"
										height="32"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										strokeWidth="1.5"
										className="mx-auto text-muziboo-text-muted/50 mb-2 group-hover:text-muziboo-gold transition-colors"
									>
										<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
										<polyline points="17 8 12 3 7 8" />
										<line x1="12" y1="3" x2="12" y2="15" />
									</svg>
									<p className="text-muziboo-text-muted group-hover:text-muziboo-text transition-colors">
										Click to select audio file
									</p>
									<p className="text-muziboo-text-muted/50 text-sm mt-1">
										MP3, WAV, FLAC, M4A, OGG · Max 50 MB
									</p>
								</div>
							)}
						</button>
					</div>

					{/* Artwork */}
					<div>
						<label className="block text-sm font-medium text-muziboo-text mb-2">
							Artwork (optional)
						</label>
						<input
							ref={artworkInputRef}
							type="file"
							accept="image/*"
							onChange={handleArtworkChange}
							className="hidden"
						/>
						<button
							type="button"
							onClick={() => artworkInputRef.current?.click()}
							className="border-2 border-dashed border-muziboo-border rounded-xl overflow-hidden hover:border-muziboo-gold/50 transition-colors group w-40 h-40"
						>
							{artworkPreview ? (
								<img
									src={artworkPreview}
									alt="Artwork preview"
									className="w-full h-full object-cover"
								/>
							) : (
								<div className="w-full h-full flex flex-col items-center justify-center">
									<svg
										width="24"
										height="24"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										strokeWidth="1.5"
										className="text-muziboo-text-muted/50 mb-1 group-hover:text-muziboo-gold transition-colors"
									>
										<rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
										<circle cx="8.5" cy="8.5" r="1.5" />
										<polyline points="21 15 16 10 5 21" />
									</svg>
									<span className="text-xs text-muziboo-text-muted/50 group-hover:text-muziboo-text-muted transition-colors">
										Add artwork
									</span>
								</div>
							)}
						</button>
					</div>

					{/* Title */}
					<div>
						<label
							htmlFor="title"
							className="block text-sm font-medium text-muziboo-text mb-1.5"
						>
							Title *
						</label>
						<input
							id="title"
							type="text"
							value={title}
							onChange={(e) => setTitle(e.target.value)}
							required
							placeholder="What's this track called?"
							className="w-full px-4 py-3 bg-muziboo-bg border border-muziboo-border rounded-lg text-muziboo-text placeholder:text-muziboo-text-muted/50 focus:outline-none focus:border-muziboo-gold transition-colors"
						/>
					</div>

					{/* Genre */}
					<div>
						<label
							htmlFor="genre"
							className="block text-sm font-medium text-muziboo-text mb-1.5"
						>
							Genre
						</label>
						<select
							id="genre"
							value={genre}
							onChange={(e) => setGenre(e.target.value)}
							className="w-full px-4 py-3 bg-muziboo-bg border border-muziboo-border rounded-lg text-muziboo-text focus:outline-none focus:border-muziboo-gold transition-colors"
						>
							<option value="">Select genre (optional)</option>
							{genres.map((g) => (
								<option key={g} value={g}>
									{g}
								</option>
							))}
						</select>
					</div>

					{/* Description */}
					<div>
						<label
							htmlFor="description"
							className="block text-sm font-medium text-muziboo-text mb-1.5"
						>
							Description
						</label>
						<textarea
							id="description"
							value={description}
							onChange={(e) => setDescription(e.target.value)}
							rows={3}
							placeholder="Tell us about this track... What's the story? What gear did you use?"
							className="w-full px-4 py-3 bg-muziboo-bg border border-muziboo-border rounded-lg text-muziboo-text placeholder:text-muziboo-text-muted/50 focus:outline-none focus:border-muziboo-gold transition-colors resize-none"
						/>
					</div>

					{error && (
						<div className="text-muziboo-red text-sm bg-muziboo-red/10 rounded-lg px-4 py-3">
							{error}
						</div>
					)}

					<button
						type="submit"
						disabled={loading || !audioFile || !title}
						className="w-full bg-muziboo-gold text-muziboo-bg font-semibold px-6 py-3.5 rounded-lg hover:bg-muziboo-text hover:text-muziboo-bg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{loading ? uploadProgress || "Uploading..." : "Upload Track"}
					</button>
				</form>
			</div>
		</div>
	);
}
