import { chromium } from "playwright";

(async () => {
	const browser = await chromium.launch();
	const context = await browser.newContext();
	const page = await context.newPage();

	page.on("console", (msg) => console.log("BROWSER LOG:", msg.text()));
	page.on("pageerror", (error) => console.log("BROWSER ERROR:", error.message));

	try {
		await page.goto("http://localhost:8080/app/login", { timeout: 10000 });
		console.log("Page loaded");
		await page.waitForTimeout(1000);

		console.log("Clicking Google button");
		await page.click("#google-btn");

		await page.waitForTimeout(2000);
		const btnText = await page.$eval("#google-btn", (el) => el.textContent);
		console.log("Button text after click:", btnText.trim());
	} catch (err) {
		console.error("SCRIPT ERROR:", err);
	}

	await browser.close();
})();
