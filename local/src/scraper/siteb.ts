import { Browser, Page } from 'puppeteer';
import { evaluateJobInterest } from '../interest';
import { JobService, Job } from '../jobService';
const urls: string[] = [];

export const scrapB = async(browser : Browser, link : string, jobCount : number, jobService: JobService) => {
    let page: Page | undefined;
    try{
        const tabs = await browser.pages();
        for (const tab of tabs) {
            if (tab.url().includes("linkedin.com")) {
                page = tab;
                break;
            }
        }
        if (!page) {
            throw new Error('Page with matching link not found');
        }

        await page.waitForSelector('ul.scaffold-layout__list-container div[data-job-id]');
        const jobDivs = document.querySelectorAll('ul.scaffold-layout__list-container div[data-job-id]');





       

    } catch (error) {
        const err = error as Error;
        let currentUrl: string = 'first page provided';
        if (urls.length) {
            currentUrl = urls[urls.length - 1];
        }
        throw new Error(`An error occurred at ${currentUrl}: ${err?.stack}`);
    }
}