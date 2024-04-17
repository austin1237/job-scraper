import fs from 'fs';

// Create a text document named with today's date
const date = new Date();

const fileName = `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}-jobs.txt`;

export function writeToFile(pageUrl: string, jobTitle: string | null, companyName: string | null, matchedKeyword: string) {
    const fileContent = `${pageUrl}, ${jobTitle}, ${companyName}, ${matchedKeyword}\n`;
    fs.appendFileSync(fileName, fileContent);
}
