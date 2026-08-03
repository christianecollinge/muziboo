import { createClient } from "@supabase/supabase-js";

const SUPABASE_URL = "https://paydmwknrjmjlswyouem.supabase.co";
const SUPABASE_ANON_KEY =
	process.env.SUPABASE_ANON_KEY ||
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InBheWRtd2tucmptamxzd3lvdWVtIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzI5NzE5MTcsImV4cCI6MjA4ODU0NzkxN30.hQes4l7_m3wER6UJ_zok0DdSNtPYt0fCMQbBTXosx7g";

const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);

async function test() {
	const { data, error } = await supabase.auth.signInWithOAuth({
		provider: "google",
		options: { redirectTo: "http://localhost:4321/app/dashboard" },
	});
	console.log("DATA:", data);
	console.log("ERROR:", error);
}

test();
