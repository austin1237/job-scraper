const puppeteer = require("puppeteer-core");
const chromium = require("@sparticuz/chromium");

const getPageHtml = async (url) => {
  const browser = await puppeteer.launch({
    args: [...chromium.args, '--disable-gpu'],
    defaultViewport: chromium.defaultViewport,
    executablePath: await chromium.executablePath(),
    headless: chromium.headless,
  });
  const page = await browser.newPage();
  await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36');
  await page.goto(url);
  const html = await page.content();
  await browser.close();
  return html;
}

exports.handler = async (event, context) => {
  const url = event.queryStringParameters.url;
  try {
    const html = await getPageHtml(url)
    const response = {
      "statusCode": 200,
      "headers": {},
      "body": html,
      "isBase64Encoded": false
    }
    return response
  } catch (error) {
    console.error(error);
    const response = {
      "statusCode": 404,
      "headers": {},
      "body": JSON.stringify(error),
      "isBase64Encoded": false
    };
    return response
  };
};