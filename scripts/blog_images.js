import fs from 'fs';
import path from 'path';

const USED_IMAGES_FILE = path.join(process.cwd(), 'scripts', 'used_images.json');
const FOLDER_ID = "1NmHmkjOQjfgbHnLixvsKLyVX6YgNFQzM";

/**
 * Get all available images from the Drive folder via gws CLI
 */
export async function getAllImages() {
    // We use a simplified version of the gws command
    const { execSync } = await import('child_process');
    try {
        const output = execSync(`gws drive files list --params '{"q": "\\"${FOLDER_ID}\\" in parents and mimeType contains \\"image/\\""}'`).toString();
        const data = JSON.parse(output);
        return data.files || [];
    } catch (e) {
        console.error("❌ Error listing Drive files:", e);
        return [];
    }
}

export function getUsedImages() {
    if (fs.existsSync(USED_IMAGES_FILE)) {
        return JSON.parse(fs.readFileSync(USED_IMAGES_FILE, 'utf-8'));
    }
    return [];
}

export function saveUsedImage(fileId) {
    const used = getUsedImages();
    if (!used.includes(fileId)) {
        used.push(fileId);
        fs.writeFileSync(USED_IMAGES_FILE, JSON.stringify(used, null, 2));
    }
}

/**
 * Returns the next available image and warns if low
 */
export async function getNextImage() {
    const allImages = await getAllImages();
    const usedIds = getUsedImages();
    
    const available = allImages.filter(img => !usedIds.includes(img.id));
    
    if (available.length === 0) {
        console.warn("⚠️ WARNING: All images in the Drive folder have been used!");
        return null;
    }
    
    if (available.length <= 3) {
        console.warn(`⚠️ WARNING: Only ${available.length} images left in the Drive folder!`);
    }
    
    // Pick a random one
    const randomIndex = Math.floor(Math.random() * available.length);
    const selected = available[randomIndex];
    
    // We don't save it here yet, we save it when it's actually applied to a post
    return selected;
}

/**
 * Injects an image into a markdown post
 */
export function injectImageIntoMarkdown(markdown, imageUrl, imageName) {
    const imageMarkdown = `\n![${imageName}](${imageUrl})\n`;
    
    // Try to inject after the H1 title
    if (markdown.includes('# ')) {
        const parts = markdown.split('\n');
        const h1Index = parts.findIndex(line => line.startsWith('# '));
        if (h1Index !== -1) {
            parts.splice(h1Index + 1, 0, imageMarkdown);
            return parts.join('\n');
        }
    }
    
    // Fallback: inject after frontmatter
    if (markdown.startsWith('---')) {
        const parts = markdown.split('---');
        if (parts.length >= 3) {
            parts[2] = imageMarkdown + parts[2];
            return parts.join('---');
        }
    }
    
    return imageMarkdown + markdown;
}

/**
 * Drive Image Link Helper (requires making images public or using a proxy)
 * For now we'll use the direct thumbnail/view link which works if public.
 */
export function getImageUrl(fileId) {
    return `https://lh3.googleusercontent.com/u/0/d/${fileId}`;
}
