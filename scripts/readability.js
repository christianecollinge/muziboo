/**
 * Simple Flesch-Kincaid Readability Utility
 */

export function countSyllables(word) {
    word = word.toLowerCase();
    if (word.length <= 3) return 1;
    word = word.replace(/(?:[^laeiouy]es|ed|[^laeiouy]e)$/, '');
    word = word.replace(/^y/, '');
    const syllables = word.match(/[aeiouy]{1,2}/g);
    return syllables ? syllables.length : 1;
}

export function calculateFleschReadingEase(text) {
    // Remove Markdown formatting
    const plainText = text
        .replace(/#+\s+.+/g, '') // Remove headers
        .replace(/---[\s\S]*?---/g, '') // Remove frontmatter
        .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1') // Remove links
        .replace(/[*_~`]+/g, '') // Remove bold/italic/code
        .trim();

    const sentences = plainText.split(/[.!?]+\s+/).filter(s => s.trim().length > 0);
    const words = plainText.split(/\s+/).filter(w => w.trim().length > 0);
    
    if (sentences.length === 0 || words.length === 0) return 0;

    let totalSyllables = 0;
    for (const word of words) {
        totalSyllables += countSyllables(word.replace(/[^a-zA-Z]/g, ''));
    }

    const avgSentenceLength = words.length / sentences.length;
    const avgSyllablesPerWord = totalSyllables / words.length;

    // Flesch Reading Ease Formula:
    // 206.835 - (1.015 * ASL) - (84.6 * ASW)
    const score = 206.835 - (1.015 * avgSentenceLength) - (84.6 * avgSyllablesPerWord);
    
    return {
        score: Math.round(score * 10) / 10,
        sentences: sentences.length,
        words: words.length,
        syllables: totalSyllables,
        avgSentenceLength: Math.round(avgSentenceLength * 10) / 10,
        avgSyllablesPerWord: Math.round(avgSyllablesPerWord * 10) / 10
    };
}

export function calculateFleschKincaidGradeLevel(text) {
    const data = calculateFleschReadingEase(text);
    if (data.words === 0) return 0;

    // Flesch–Kincaid Grade Level Formula:
    // 0.39 * (words / sentences) + 11.8 * (syllables / words) - 15.59
    const gradeLevel = (0.39 * data.avgSentenceLength) + (11.8 * data.avgSyllablesPerWord) - 15.59;
    
    return Math.round(gradeLevel * 10) / 10;
}
