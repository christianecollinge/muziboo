import 'dotenv/config';
import { GoogleGenAI } from '@google/genai';
import fs from 'fs';
import path from 'path';

const ai = new GoogleGenAI({ apiKey: process.env.GEMINI_API_KEY });

async function dfsPost(endpoint, body) {
    const credentials = Buffer.from(`${process.env.DATAFORSEO_LOGIN}:${process.env.DATAFORSEO_PASSWORD}`).toString('base64');
    const res = await fetch(`https://api.dataforseo.com${endpoint}`, {
        method: 'POST',
        headers: { 'Authorization': `Basic ${credentials}`, 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    });
    return res.json();
}

const BLOG_DIR = path.join(process.cwd(), 'src', 'content', 'blog');

const TARGETS = [
    { location: "United States", language: "English", instructions: "Write the post in English, focusing on the US music industry." },
    { location: "United Kingdom", language: "English", instructions: "Write the post in British English, focusing on the UK music scene." },
    { location: "Spain", language: "Spanish", instructions: "Write the post entirely in Spanish, focusing on the Spanish music scene." },
    { location: "Germany", language: "German", instructions: "Write the post entirely in German, focusing on the German music scene." }
];

const dayOfYear = Math.floor((new Date() - new Date(new Date().getFullYear(), 0, 0)) / (1000 * 60 * 60 * 24));
const currentTarget = TARGETS[dayOfYear % TARGETS.length];

async function fetchKeywords(target) {
    console.log(`🔍 Fetching keywords from DataForSEO for ${target.location} (${target.language})...`);
    try {
        const seedKeywords = [
            "upload original music",
            "get feedback on unfinished songs",
            "where to share bedroom pop demos",
            "community for hobby musicians",
            "indie musician feedback forum"
        ];

        const post_array = seedKeywords.map(kw => ({
            "location_name": target.location,
            "language_name": target.language,
            "keyword": kw,
            "limit": 2
        }));

        const response = await dfsPost("/v3/dataforseo_labs/google/keyword_suggestions/live", post_array);

        if (response.tasks && response.tasks.length > 0) {
            const task = response.tasks[0];
            if (task.result && task.result.length > 0 && task.result[0].items) {
                const items = task.result[0].items || [];
                return items.slice(0, 10).map(item => item.keyword);
            } else {
                console.error("❌ DataForSEO returned an error or empty result payload:", JSON.stringify(task, null, 2));
                throw new Error("DataForSEO API failed to return keywords");
            }
        } else {
            console.error("❌ DataForSEO response format unexpected:", JSON.stringify(response, null, 2));
            throw new Error("DataForSEO API failed to return tasks");
        }
    } catch (error) {
        console.error("❌ DataForSEO Fetch/Parse Error:", error);
        throw error; // Fail loudly instead of using static fallbacks indefinitely
    }
}

async function generateBlogPost(keyword, target) {
    console.log(`✍️ Generating blog post for keyword: "${keyword}" in ${target.language}...`);

    const prompt = `
    You are an expert music industry blogger writing for 'Muziboo'.

    ${target.instructions}

    Muziboo is a workshop for real music and real people—not a streaming platform or a talent show. It is a space for hobby musicians, bedroom producers, and anyone who makes music with their own hands. Users share demos, fragments, late-night recordings, and unfinished songs to get constructive feedback from people who care about craft, not metrics, likes, or ranking games.

    Write an SEO-optimized, engaging blog post targeting the core keyword: "${keyword}".

    IMPORTANT FORMATTING RULES:
    - Do NOT use bullet points or numbered lists anywhere in the post.
    - Write in flowing prose paragraphs only.
    - Use H2 and H3 subheadings to break up sections.

    Output the response STRICTLY as a Markdown file with YAML frontmatter. Do not include any text outside of the Markdown block.

    Structure the markdown exactly like this:
    ---
    title: "Catchy SEO title including the keyword"
    description: "A short 160-character SEO description."
    pubDate: ${new Date().toISOString()}
    author: "Muziboo Team"
    tags: ["music", "creators", "community"]
    ---
    # Your catchy H1 title here

    [Rest of the blog post with flowing prose paragraphs and subheadings — NO bullet points]

    Make sure to mention Muziboo respectfully throughout the post as the best workshop to upload unpolished music and get real human feedback. The content should be at least 600 words long.
    `;

    try {
        const response = await ai.models.generateContent({
            model: 'gemini-2.5-flash',
            contents: prompt,
        });

        let markdown = response.text;

        if (markdown.startsWith('```markdown')) {
            markdown = markdown.replace(/^```markdown\n/, '').replace(/\n```$/, '');
        } else if (markdown.startsWith('```')) {
            markdown = markdown.replace(/^```\n/, '').replace(/\n```$/, '');
        }

        return markdown;
    } catch (error) {
        console.error(`❌ Gemini Error for "${keyword}":`, error);
        return null;
    }
}

async function main() {
    if (!process.env.DATAFORSEO_LOGIN || !process.env.DATAFORSEO_PASSWORD || !process.env.GEMINI_API_KEY) {
        console.error("❌ Missing environment variables. Please check your .env file.");
        process.exit(1);
    }

    if (!fs.existsSync(BLOG_DIR)) fs.mkdirSync(BLOG_DIR, { recursive: true });

    const keywords = await fetchKeywords(currentTarget);
    console.log(`✅ Found ${keywords.length} keywords.`);

    for (const keyword of keywords) {
        const slug = keyword.toLowerCase().replace(/[^a-z0-9]+/g, '-');
        const markdown = await generateBlogPost(keyword, currentTarget);
        if (markdown) {
            const filepath = path.join(BLOG_DIR, `${slug}.md`);
            fs.writeFileSync(filepath, markdown);
            console.log(`✅ Saved: ${slug}.md`);
        }

        // Wait 4 seconds between posts to respect rate limits
        await new Promise(resolve => setTimeout(resolve, 4000));
    }

    console.log("🎉 All blog posts generated successfully!");
}

main();
