import 'dotenv/config';
import { GoogleGenAI } from '@google/genai';
import fs from 'fs';
import path from 'path';
import { calculateFleschReadingEase, calculateFleschKincaidGradeLevel } from './readability.js';
import { getRelevantImage, downloadImage, saveUsedImage } from './blog_images.js';

const ai = new GoogleGenAI({ apiKey: process.env.GEMINI_API_KEY });

async function dfsPost(endpoint, body) {
    const credentials = Buffer.from(`${process.env.DATAFORSEO_LOGIN}:${process.env.DATAFORSEO_PASSWORD}`).toString('base64');
    const res = await fetch(`https://api.dataforseo.com${endpoint}`, {
        method: 'POST',
        headers: { 'Authorization': `Basic ${credentials}`, 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    });
    if (!res.ok) {
        const text = await res.text();
        throw new Error(`DataForSEO API error (${res.status}): ${text}`);
    }
    return res.json();
}

const BLOG_DIR = path.join(process.cwd(), 'src', 'content', 'blog');

const TARGETS = [
    { location: "United States", language: "English", instructions: "Write the post in English, focusing on the US music industry." },
    { location: "United Kingdom", language: "English", instructions: "Write the post in British English, focusing on the UK music scene." },
    { location: "Spain", language: "Spanish", instructions: "Write the post entirely in Spanish, focusing on the Spanish music scene." },
    { location: "Germany", language: "German", instructions: "Write the post entirely in German, focusing on the German music scene." }
];

const FALLBACK_KEYWORDS = {
    "English": [
        "best community for musicians 2026",
        "how to share music demos for feedback",
        "where to post unfinished songs",
        "bedroom producer feedback groups",
        "human feedback on music without algorithms"
    ],
    "Spanish": [
        "mejor comunidad para músicos 2026",
        "cómo compartir maquetas de música para recibir feedback",
        "donde subir canciones sin terminar",
        "grupos de feedback para productores de dormitorio",
        "feedback humano sobre música sin algoritmos"
    ],
    "German": [
        "beste Community für Musiker 2026",
        "Musik-Demos für Feedback teilen",
        "wo man unfertige Songs hochlädt",
        "Feedback-Gruppen für Bedroom Producer",
        "menschliches Feedback zu Musik ohne Algorithmen"
    ]
};

const dayOfYear = Math.floor((new Date() - new Date(new Date().getFullYear(), 0, 0)) / (1000 * 60 * 60 * 24));
const currentTarget = TARGETS[dayOfYear % TARGETS.length];

async function fetchKeywords(target) {
    console.log(`🔍 Fetching keywords from DataForSEO for ${target.location} (${target.language})...`);
    
    if (!process.env.DATAFORSEO_LOGIN || !process.env.DATAFORSEO_PASSWORD) {
        console.warn(`⚠️ DataForSEO credentials are missing. Skipping API and falling back to static keywords.`);
        return (FALLBACK_KEYWORDS[target.language] || FALLBACK_KEYWORDS["English"]).slice(0, 1);
    }

    try {
        const seedKeywords = [
            "upload original music",
            "get feedback on unfinished songs"
        ];

        const post_array = seedKeywords.map(kw => ({
            "location_name": target.location,
            "language_name": target.language,
            "keyword": kw,
            "limit": 1
        }));

        const response = await dfsPost("/v3/dataforseo_labs/google/keyword_suggestions/live", post_array);

        if (response.tasks && response.tasks.length > 0) {
            const task = response.tasks[0];
            if (task.result && task.result.length > 0 && task.result[0].items && task.result[0].items.length > 0) {
                const items = task.result[0].items || [];
                // Only take 1 keyword to stay under free tier daily limits
                return items.slice(0, 1).map(item => item.keyword);
            }
        }
        
        console.warn(`⚠️ DataForSEO returned no results. Using 1 fallback keyword.`);
        return (FALLBACK_KEYWORDS[target.language] || FALLBACK_KEYWORDS["English"]).slice(0, 1);

    } catch (error) {
        console.error("❌ DataForSEO Fetch/Parse Error:", error);
        return (FALLBACK_KEYWORDS[target.language] || FALLBACK_KEYWORDS["English"]).slice(0, 1);
    }
}

async function generateBlogPost(keyword, target, pubDate, retryCount = 0) {
    console.log(`✍️ Generating blog post for keyword: "${keyword}"... (Attempt ${retryCount + 1})`);

    const prompt = `
    Write a 500-word SEO blog post for 'Muziboo' about: "${keyword}".
    ${target.instructions}

    Muziboo is a workshop for real music and real people to share demos and get human feedback. 

    RULES:
    - Flesch Reading Ease: 60+ (Simple language).
    - Format: Prose paragraphs with H2/H3 subheadings.
    - NO bullet points or lists.
    - Mention Muziboo as the best place for human feedback on unpolished music.

    Output as Markdown with YAML frontmatter (title, description, pubDate: ${pubDate}, author: "Muziboo Team", tags: ["music", "creators"]).
    `;

    try {
        const response = await ai.models.generateContent({
            model: 'models/gemini-2.0-flash',
            contents: prompt,
        });

        let markdown = response.text;
        if (!markdown) throw new Error("Gemini returned empty response");

        // Clean up markdown block if present
        markdown = markdown.trim();
        if (markdown.includes('```markdown')) {
            markdown = markdown.match(/```markdown\n([\s\S]*?)\n```/)?.[1] || markdown;
        } else if (markdown.includes('```')) {
            markdown = markdown.match(/```\n?([\s\S]*?)\n?```/)?.[1] || markdown;
        }

        return markdown.trim();
    } catch (error) {
        if ((error.status === 429 || error.message?.includes('quota')) && retryCount < 3) {
            const waitTime = 300000; // 5 minutes
            console.warn(`⚠️ Rate limit hit for "${keyword}". Waiting 5 minutes before retry...`);
            await new Promise(resolve => setTimeout(resolve, waitTime));
            return generateBlogPost(keyword, target, pubDate, retryCount + 1);
        }
        console.error(`❌ Gemini Error for "${keyword}":`, error);
        return null;
    }
}

async function main() {
    const missingVars = [];
    if (!process.env.GEMINI_API_KEY) missingVars.push("GEMINI_API_KEY");

    if (missingVars.length > 0) {
        console.error(`❌ Missing environment variables: ${missingVars.join(", ")}. Please check your .env file or repository secrets.`);
        process.exit(1);
    }

    if (!fs.existsSync(BLOG_DIR)) fs.mkdirSync(BLOG_DIR, { recursive: true });

    const keywords = await fetchKeywords(currentTarget);
    console.log(`✅ Found ${keywords.length} keywords.`);

    let i = 0;
    for (const keyword of keywords) {
        const slug = keyword.toLowerCase().replace(/[^a-z0-9]+/g, '-');
        
        let pubDate;
        let logDateString;
        if (i < 2) {
            // Keep 2 posts exactly on today's date
            pubDate = new Date().toISOString();
            logDateString = "today";
        } else {
            // SEO: Backdate randomly between 1 and 180 days ago to make historical publishing look organic
            const randomDaysAgo = Math.floor(Math.random() * 180) + 1;
            pubDate = new Date(Date.now() - randomDaysAgo * 24 * 60 * 60 * 1000).toISOString();
            logDateString = `backdated ${randomDaysAgo} days`;
        }
        i++;
        
        let markdown = await generateBlogPost(keyword, currentTarget, pubDate);
        if (markdown) {
            const { image, alt } = await getRelevantImage(markdown, ai);
            if (image) {
                const localUrl = await downloadImage(image.id, image.name);
                if (localUrl) {
                    // Insert local image path into frontmatter
                    markdown = markdown.replace('---', `---\nimage: "${localUrl}"`);
                    saveUsedImage(image.id);
                    console.log(`🖼️ Downloaded relevant image: ${image.name} to ${localUrl}`);
                }
            }

            const filepath = path.join(BLOG_DIR, `${slug}.md`);
            fs.writeFileSync(filepath, markdown);
            
            const readability = calculateFleschReadingEase(markdown);
            const gradeLevel = calculateFleschKincaidGradeLevel(markdown);
            
            console.log(`✅ Saved: ${slug}.md (${logDateString})`);
            console.log(`   Readability: ${readability.score} (Ease), Grade Level: ${gradeLevel}`);
        }

        // Wait 15 seconds between posts to respect rate limits
        await new Promise(resolve => setTimeout(resolve, 15000));
    }

    console.log("🎉 All blog posts generated successfully!");
}

main().catch(error => {
    console.error("🔥 FATAL ERROR DURING EXECUTION:");
    console.error(error);
    if (error.stack) console.error(error.stack);
    process.exit(1);
});
