const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();
  
  page.on('console', msg => console.log('BROWSER LOG:', msg.text()));
  page.on('pageerror', error => console.log('BROWSER ERROR:', error.message));

  try {
    await page.goto('http://localhost:4321/app/login', { timeout: 10000 });
    console.log('Page loaded');
    await page.waitForTimeout(2000);
    
    console.log('Clicking Google button');
    await page.click('#google-btn');
    
    await page.waitForTimeout(2000);
    console.log('Done waiting');
  } catch (err) {
    console.error('SCRIPT ERROR:', err);
  }
  
  await browser.close();
})();
