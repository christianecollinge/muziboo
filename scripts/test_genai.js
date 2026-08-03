import "dotenv/config";
import { GoogleGenAI } from "@google/genai";

async function test() {
	const apiKey = process.env.GEMINI_API_KEY;
	if (!apiKey) {
		console.error("No API key found");
		return;
	}
	const ai = new GoogleGenAI({ apiKey });
	try {
		console.log("Generating content...");
		const response = await ai.models.generateContent({
			model: "models/gemini-1.5-flash",
			contents: "Say hello",
		});
		console.log("Response text:", response.text);
	} catch (e) {
		console.error("Error:", e);
	}
}

test();
