import React, { useState } from "react";
import { supabase } from "../lib/supabase";

interface AuthFormProps {
	mode: "login" | "signup";
}

export default function AuthForm({ mode }: AuthFormProps) {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [username, setUsername] = useState("");
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError("");
		setLoading(true);

		try {
			if (mode === "signup") {
				const { error: signUpError } = await supabase.auth.signUp({
					email,
					password,
					options: {
						data: {
							username: username.toLowerCase().replace(/[^a-z0-9_-]/g, ""),
							display_name: username,
						},
					},
				});

				if (signUpError) throw signUpError;

				window.location.href = "/dashboard";
			} else {
				const { error: signInError } = await supabase.auth.signInWithPassword({
					email,
					password,
				});

				if (signInError) throw signInError;

				window.location.href = "/dashboard";
			}
		} catch (err: any) {
			setError(err.message || "An error occurred");
		} finally {
			setLoading(false);
		}
	};

	return (
		<div className="w-full max-w-md mx-auto">
			<div className="bg-muziboo-bg-softer rounded-xl p-8 border border-muziboo-border">
				<h1 className="text-3xl font-semibold text-muziboo-text mb-2">
					{mode === "login" ? "Welcome back" : "Join Muziboo"}
				</h1>
				<p className="text-muziboo-text-muted mb-8">
					{mode === "login"
						? "Log in to your workshop"
						: "Create your space for real music"}
				</p>

				<form onSubmit={handleSubmit} className="space-y-5">
					{mode === "signup" && (
						<div>
							<label
								htmlFor="username"
								className="block text-sm font-medium text-muziboo-text mb-1.5"
							>
								Username
							</label>
							<input
								id="username"
								type="text"
								value={username}
								onChange={(e) => setUsername(e.target.value)}
								required
								minLength={3}
								maxLength={30}
								placeholder="your-username"
								className="w-full px-4 py-3 bg-muziboo-bg border border-muziboo-border rounded-lg text-muziboo-text placeholder:text-muziboo-text-muted/50 focus:outline-none focus:border-muziboo-gold transition-colors"
							/>
						</div>
					)}

					<div>
						<label
							htmlFor="email"
							className="block text-sm font-medium text-muziboo-text mb-1.5"
						>
							Email
						</label>
						<input
							id="email"
							type="email"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							required
							placeholder="you@example.com"
							className="w-full px-4 py-3 bg-muziboo-bg border border-muziboo-border rounded-lg text-muziboo-text placeholder:text-muziboo-text-muted/50 focus:outline-none focus:border-muziboo-gold transition-colors"
						/>
					</div>

					<div>
						<label
							htmlFor="password"
							className="block text-sm font-medium text-muziboo-text mb-1.5"
						>
							Password
						</label>
						<input
							id="password"
							type="password"
							value={password}
							onChange={(e) => setPassword(e.target.value)}
							required
							minLength={6}
							placeholder="At least 6 characters"
							className="w-full px-4 py-3 bg-muziboo-bg border border-muziboo-border rounded-lg text-muziboo-text placeholder:text-muziboo-text-muted/50 focus:outline-none focus:border-muziboo-gold transition-colors"
						/>
					</div>

					{error && (
						<div className="text-muziboo-red text-sm bg-muziboo-red/10 rounded-lg px-4 py-3">
							{error}
						</div>
					)}

					<button
						type="submit"
						disabled={loading}
						className="w-full bg-muziboo-gold text-muziboo-bg font-semibold px-6 py-3.5 rounded-lg hover:bg-muziboo-text hover:text-muziboo-bg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{loading
							? "Loading..."
							: mode === "login"
								? "Log In"
								: "Create Account"}
					</button>
				</form>

				<p className="mt-6 text-center text-muziboo-text-muted text-sm">
					{mode === "login" ? (
						<>
							Don't have an account?{" "}
							<a
								href="/signup"
								className="text-muziboo-gold hover:text-muziboo-text transition-colors"
							>
								Sign up
							</a>
						</>
					) : (
						<>
							Already have an account?{" "}
							<a
								href="/login"
								className="text-muziboo-gold hover:text-muziboo-text transition-colors"
							>
								Log in
							</a>
						</>
					)}
				</p>
			</div>
		</div>
	);
}
