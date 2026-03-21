import { createClient } from "@supabase/supabase-js";

const supabaseUrl =
	import.meta.env.PUBLIC_SUPABASE_URL ||
	"https://yndossoknynhepphdphd.supabase.co";
// Provide a dummy key if env var is missing to prevent createClient from throwing an error and crashing the whole React app
const supabaseAnonKey =
	import.meta.env.PUBLIC_SUPABASE_ANON_KEY ||
	"dummy_key_to_prevent_crash_please_set_env_var";

export const supabase = createClient(supabaseUrl, supabaseAnonKey);

/**
 * Get the current session access token (JWT).
 * Returns null if not authenticated.
 */
export async function getAccessToken(): Promise<string | null> {
	const {
		data: { session },
	} = await supabase.auth.getSession();
	return session?.access_token ?? null;
}

/**
 * Get the current user's ID.
 */
export async function getCurrentUserId(): Promise<string | null> {
	const {
		data: { session },
	} = await supabase.auth.getSession();
	return session?.user?.id ?? null;
}
