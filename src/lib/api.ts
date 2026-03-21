const API_URL = import.meta.env.PUBLIC_API_URL || "http://localhost:8080";

interface FetchOptions {
	method?: string;
	body?: any;
	token?: string | null;
	isFormData?: boolean;
}

/**
 * Make an API call to the Go backend.
 */
export async function api<T = any>(
	path: string,
	options: FetchOptions = {}
): Promise<T> {
	const { method = "GET", body, token, isFormData = false } = options;

	const headers: Record<string, string> = {};

	if (token) {
		headers["Authorization"] = `Bearer ${token}`;
	}

	if (!isFormData) {
		headers["Content-Type"] = "application/json";
	}

	const fetchOptions: RequestInit = {
		method,
		headers,
	};

	if (body) {
		fetchOptions.body = isFormData ? body : JSON.stringify(body);
	}

	const response = await fetch(`${API_URL}${path}`, fetchOptions);

	if (!response.ok) {
		const error = await response
			.json()
			.catch(() => ({ error: "Request failed" }));
		throw new Error(error.error || `API error: ${response.status}`);
	}

	if (response.status === 204) {
		return null as T;
	}

	return response.json();
}

/**
 * Upload a file to the Go backend.
 */
export async function uploadFile(
	endpoint: "/api/upload/audio" | "/api/upload/artwork",
	file: File,
	token: string,
	queryParams?: string
): Promise<{ url: string }> {
	const formData = new FormData();
	formData.append("file", file);

	const path = queryParams ? `${endpoint}?${queryParams}` : endpoint;

	return api(path, {
		method: "POST",
		body: formData,
		token,
		isFormData: true,
	});
}
