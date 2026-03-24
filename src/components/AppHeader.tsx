import React, { useState, useEffect } from "react";
import { supabase } from "../lib/supabase";

interface AppHeaderProps {
	currentPage?: string;
}

export default function AppHeader({ currentPage }: AppHeaderProps) {
	const [user, setUser] = useState<any>(null);
	const [username, setUsername] = useState("");
	const [menuOpen, setMenuOpen] = useState(false);

	useEffect(() => {
		supabase.auth.getSession().then(({ data: { session } }) => {
			setUser(session?.user || null);
			if (session?.user) {
				setUsername(session.user.user_metadata?.username || "");
			}
		});

		const {
			data: { subscription },
		} = supabase.auth.onAuthStateChange((_event, session) => {
			setUser(session?.user || null);
			if (session?.user) {
				setUsername(session.user.user_metadata?.username || "");
			}
		});

		return () => subscription.unsubscribe();
	}, []);

	const handleLogout = async () => {
		await supabase.auth.signOut();
		window.location.href = "/";
	};

	const navLinks = [{ href: "/app/explore", label: "Explore" }];

	const isActive = (href: string) => currentPage === href;

	return (
		<header className="sticky top-0 z-50 bg-muziboo-bg/90 backdrop-blur-md border-b border-muziboo-border">
			<div className="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
				{/* Logo */}
				<a
					href="/"
					className="text-xl font-bold text-muziboo-gold tracking-tight hover:text-muziboo-text transition-colors"
				>
					muziboo
				</a>

				{/* Nav */}
				<nav className="flex items-center gap-6">
					{navLinks.map((link) => (
						<a
							key={link.href}
							href={link.href}
							className={`text-sm font-medium transition-colors ${
								isActive(link.href)
									? "text-muziboo-gold"
									: "text-muziboo-text-muted hover:text-muziboo-text"
							}`}
						>
							{link.label}
						</a>
					))}

					{user ? (
						<>
							<a
								href="/app/upload"
								className={`text-sm font-medium transition-colors ${
									isActive("/upload")
										? "text-muziboo-gold"
										: "text-muziboo-text-muted hover:text-muziboo-text"
								}`}
							>
								Upload
							</a>

							{/* User menu */}
							<div className="relative">
								<button
									onClick={() => setMenuOpen(!menuOpen)}
									className="flex items-center gap-2 text-sm font-medium text-muziboo-text-muted hover:text-muziboo-text transition-colors"
								>
									<div className="w-8 h-8 rounded-full bg-muziboo-gold/20 flex items-center justify-center">
										<span className="text-xs text-muziboo-gold font-semibold">
											{(username || user.email || "?").charAt(0).toUpperCase()}
										</span>
									</div>
									<svg
										width="12"
										height="12"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										strokeWidth="2"
										className={`transition-transform ${menuOpen ? "rotate-180" : ""}`}
									>
										<polyline points="6 9 12 15 18 9" />
									</svg>
								</button>

								{menuOpen && (
									<>
										<div
											className="fixed inset-0"
											onClick={() => setMenuOpen(false)}
										/>
										<div className="absolute right-0 top-12 w-48 bg-muziboo-bg-softer border border-muziboo-border rounded-xl shadow-xl py-2 z-50">
											<a
												href="/app/dashboard"
												className="block px-4 py-2.5 text-sm text-muziboo-text hover:bg-muziboo-gold/10 transition-colors"
											>
												Dashboard
											</a>
											{username && (
												<a
													href={`/user/${username}`}
													className="block px-4 py-2.5 text-sm text-muziboo-text hover:bg-muziboo-gold/10 transition-colors"
												>
													My Profile
												</a>
											)}
											<hr className="my-1 border-muziboo-border" />
											<button
												onClick={handleLogout}
												className="w-full text-left px-4 py-2.5 text-sm text-muziboo-red hover:bg-muziboo-red/10 transition-colors"
											>
												Log Out
											</button>
										</div>
									</>
								)}
							</div>
						</>
					) : null}
				</nav>
			</div>
		</header>
	);
}
