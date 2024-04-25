import puppeteer from 'puppeteer';
import { scrap } from './scraper/sitea';
import { JobService } from './jobService';

const main = async () =>{
    let jobCount = 0;
    const link = process.env.LOCAL_SCRAPER_SITEA;
    const jobApiEndpoint = process.env.JOB_API_ENDPOINT;
    if (!link) {
        console.error('LOCAL_SCRAPER_SITEA environment variable isnt set');
        return;
    }
    if (!jobApiEndpoint) {
        console.error('JOB_API_ENDPOINT environment variable isnt set');
        return;
    }
    const browser = await puppeteer.launch({ headless: false });
    const jobService = new JobService(jobApiEndpoint);
    await scrap(browser, link, jobCount, jobService);
    console.log('Scraping complete');
    await browser.close();
}

main().catch(error => {
    console.error(error);
});