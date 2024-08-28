import puppeteer from 'puppeteer';
import { scrap } from './scraper/sitea';
import { JobService } from './jobService';

const main = async () =>{
    let jobCount = 0;
    const link = process.env.LOCAL_SCRAPER_SITEA;
    const jobApiEndpoint = process.env.JOB_API_ENDPOINT;
    
    // Chrome needs to be launched with remote debugging port enabled
    const browserURL = "http://localhost:9222"

    if (!link) {
        console.error('LOCAL_SCRAPER_SITEA environment variable isnt set');
        return;
    }
    if (!jobApiEndpoint) {
        console.error('JOB_API_ENDPOINT environment variable isnt set');
        return;
    }
    const connectOptions = {
        browserURL: browserURL,
        protocolTimeout:  60000, // 60 seconds
    }
    const browser = await puppeteer.connect(connectOptions);
    const jobService = new JobService(jobApiEndpoint);
    await scrap(browser, link, jobCount, jobService);
    console.log('Scraping complete');
    await browser.disconnect();
}

main().catch(error => {
    console.error(error);
});