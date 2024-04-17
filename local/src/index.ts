import puppeteer from 'puppeteer';
import { scrap } from './scraper/sitea';

const main = async () =>{
    let jobCount = 0;
    const link = process.env.LOCAL_SCRAPER_SITEA;
    if (!link) {
        console.error('LOCAL_SCRAPER_SITEA environment variable isnt set');
        return;
    }
    const browser = await puppeteer.launch({ headless: false });
    await scrap(browser, link, jobCount);
    console.log('Scraping complete');
    await browser.close();
}

main().catch(error => {
    console.error(error);
});