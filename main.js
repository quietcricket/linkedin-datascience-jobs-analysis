const {
	LinkedinScraper,
	relevanceFilter,
	timeFilter,
	typeFilter,
	experienceLevelFilter,
	events,
} = require("linkedin-jobs-scraper");

const fs = require('fs');
(async () => {
	// Each scraper instance is associated with one browser.
	// Concurrent queries will run on different pages within the same browser instance.
	const scraper = new LinkedinScraper({
		lswMo: 500,
		args: [
			"--lang=en-US",
			'--window-size=1920,3828',
		],
	});

	// Add listeners for scraper events
	scraper.on(events.scraper.data, (data) => {

		try {
			delete data.descriptionHTML;
			let filename = `data/${data.jobId}.json`;
			if (fs.existsSync(filename)) {
				console.log('Already scraped: ' + data.jobId);
			} else {
				fs.writeFileSync(filename, JSON.stringify(data, null, 2));
			}
		} catch (err) {
			console.log(err);
		}
	});

	scraper.on(events.scraper.error, (err) => {
		console.error(err);
	});

	scraper.on(events.scraper.end, () => {
		console.log('All done!');
	});

	// Add listeners for puppeteer browser events
	scraper.on(events.puppeteer.browser.targetcreated, () => {
	});
	scraper.on(events.puppeteer.browser.targetchanged, () => {
	});
	scraper.on(events.puppeteer.browser.targetdestroyed, () => {
	});
	scraper.on(events.puppeteer.browser.disconnected, () => {
	});

	// Custom function executed on browser side to extract job description [optional]
	const descriptionFn = () => document.querySelector(".description__text")
		.innerText
		.replace(/[\s\n\r]+/g, " ")
		.trim();

	// Run queries concurrently    
	await Promise.all([
		// Run queries serially
		scraper.run([
			{
				query: "data scientist",
				options: {
					locations: ["Singapore"],
					limit: 50,
				}
			}
		])
	]);

	// Close browser
	await scraper.close();
})();