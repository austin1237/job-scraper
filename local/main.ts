import puppeteer from "puppeteer";
import { scrap } from "./scraper/sitea.ts";
import { JobService } from "./jobService.ts";

const main = async () => {
  const link = Deno.env.get("LOCAL_SCRAPER_SITEA");
  const jobApiEndpoint = Deno.env.get("JOB_API_ENDPOINT");

  // Chrome needs to be launched with remote debugging port enabled
  const browserURL = "http://localhost:9222";

  if (!link) {
    console.error("LOCAL_SCRAPER_SITEA environment variable isnt set");
    return;
  }
  if (!jobApiEndpoint) {
    console.error("JOB_API_ENDPOINT environment variable isnt set");
    return;
  }
  const connectOptions = {
    browserURL: browserURL,
    protocolTimeout: 60000, // 60 seconds
  };
  const browser = await puppeteer.connect(connectOptions);
  const jobService = new JobService(jobApiEndpoint);
  await scrap(browser, link, 0, jobService);
  console.log("Scraping complete");
  await browser.disconnect();
};

main().catch((error) => {
  console.error(error);
});
