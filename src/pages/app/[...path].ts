import type { APIRoute } from "astro";

export const ALL: APIRoute = async ({ request }) => {
	const url = new URL(request.url);
	const targetUrl = new URL(
		url.pathname + url.search,
		"https://muziboo-6n8i.onrender.com"
	);

	const proxyRequest = new Request(targetUrl, request);
	proxyRequest.headers.set("Host", "muziboo-6n8i.onrender.com");
	proxyRequest.headers.set("X-Forwarded-Host", url.host);

	return fetch(proxyRequest);
};
