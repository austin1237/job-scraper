import { Browser, Page } from 'puppeteer';
import { evaluateJobInterest } from '../interest';
import { writeToFile } from '../fileHandler';
const urls: string[] = [];

const didUrlChange = async (page: Page): Promise<boolean> => {

    const url: string = page.url();
    if (urls.includes(url)) {
        // sleep for one second and check again
        console.log('page url hasn\'t changed, sleeping for 3 second');
        await new Promise(resolve => setTimeout(resolve, 3000));
        return didUrlChange(page);
    }

    urls.push(url);
    return true;
};


export const scrap = async(browser : Browser, link : string, jobCount : number) => {
    try{    
        const page = await browser.newPage();
        await page.goto(link);
        await page.waitForSelector('.slider_container');

        while (true) {
            const containers = await page.$$('.slider_container');
            for (const container of containers) {
                const h2 = await container.$('h2');
                if (h2) {
                    await h2.click();
                    const className = await page.evaluate(element => element.className, h2);
                    if (className.startsWith('jobTitle')) {

                        // Waiting for the new job to load on the subpage
                        await page.waitForSelector('.slider_container', { visible: true });
                        await page.waitForSelector('#jobDescriptionText');

                        // subpage should have loaded by now confirm by checking the url
                        await didUrlChange(page);
                       
                        const jobTitle = await page.evaluate(element => element.textContent, h2); // Get the h2 text
                        let companyName = null;

                        // company isn't always set for w/e reason
                        try {
                            companyName = await container.$eval('[data-testid="company-name"]', element => element.textContent); // Get the company name
                        } catch (error) {
                            companyName = 'noCompanyFound';
                        }

                        const jobDescriptionText = await page.$eval('#jobDescriptionText', element => element.textContent);
                        const pageUrl = urls[urls.length - 1]; // Get the current page URL
                        const jobCategory = evaluateJobInterest(jobTitle, companyName, jobDescriptionText);

                        if(jobCategory){
                            writeToFile(pageUrl, jobTitle, companyName, jobCategory);
                            jobCount++;
                            console.log(`Job found ${jobCount}`);
                        }
                        
                        await new Promise(resolve => setTimeout(resolve, 5000)); // Sleep for 5 seconds to avoid bot detection
                    }
                }
            }
        
            // Click on the "Next" button and wait for the next page to load
            const nextButton = await page.$('[data-testid="pagination-page-next"]');
            if (nextButton) {
                await Promise.all([
                    nextButton.click(),
                    page.waitForNavigation({ waitUntil: 'networkidle0' }),
                ]);
            } else {
                break; // Exit the loop if there is no "Next" button
            }
        }
    } catch (error) {
        const err = error as Error;
        let currentUrl: string = 'first page provided';
        if (urls.length) {
            currentUrl = urls[urls.length - 1];
        }
        throw new Error(`An error occurred at ${currentUrl}: ${err?.stack}`);
    }
}