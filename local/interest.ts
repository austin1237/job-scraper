const includedKeywords = [
  "node",
  "nodejs",
  "node.js",
  "go",
  "golang",
  "typescript",
];
const excludedKeywords = ["contract", "web3", "blockchain", "crypto"];
const includeTitles = [
  "software engineer",
  "developer",
  "backend engineer",
  "backend developer",
  "backend",
  "software developer",
  "api",
];
const excludeTitles = [
  "front-end",
  "front end",
  "frontend",
  ".net",
  "java",
  "manager",
  "lead",
  "staff",
  "principal",
  "contract",
  "c#",
  "microsoft",
];
const includeRegex = new RegExp(includeTitles.join("|"), "i");
const excludeRegex = new RegExp(excludeTitles.join("|"), "i");
const excludeCompanyRegex = /(consulting|recruiting)/i;

export function evaluateJobInterest(
  jobTitle: string | null,
  companyName: string | null,
  jobDescriptionText: string | null,
): string | null {
  let matchedKeyword = null;
  let matchedExcludedKeyword = null;

  if (!companyName || !jobTitle || !jobDescriptionText) {
    return null;
  }

  for (const keyword of includedKeywords) {
    const regex = new RegExp("\\b" + keyword + "\\b", "i");
    if (jobDescriptionText && regex.test(jobDescriptionText)) {
      matchedKeyword = keyword;
      break;
    }
  }

  if (!matchedKeyword) {
    return null;
  }

  for (const keyword of excludedKeywords) {
    const regex = new RegExp("\\b" + keyword + "\\b", "i");
    if (jobDescriptionText && regex.test(jobDescriptionText)) {
      matchedExcludedKeyword = keyword;
      break;
    }
  }

  if (matchedExcludedKeyword) {
    return null;
  }

  if (
    !includeRegex.test(jobTitle) || excludeRegex.test(jobTitle) ||
    excludeCompanyRegex.test(companyName)
  ) {
    return null;
  }

  return matchedKeyword;
}
