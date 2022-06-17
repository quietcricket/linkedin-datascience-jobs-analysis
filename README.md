# LinkedIn Data Scientist Jobs Analysis


## Why

This project is planned for the final project for Data Science specialization course on Coursera. It will be fun to analyze data science with data science knowledge. 


## Getting the data

The only viable method to get job listings from LinkedIn is to use web scraping. LinkedIn does not provide a public API to query job listings.

To avoid being blocked by LinkedIn, the scraping process is split into two steps. The first scraper opens the search page and triggers the "Load More" button repeatedly until the serve stops sending back new data. The scraper extracts links to the individual jobs' detail URL and saves into a text file. After that, the second scraper fetches the actual job details page and extracts the useful data from the HTML pages.

For each search, the server returns 900 to 1000 jobs. The scraping process is repeated everyday to collect new listings. 

The keyword used is "data scientists" and the locations used are "Singapore", "United States", "United Kingdom", "Canada", "Australia" and "China".

## About the dataset

The dataset consits of [XXXXX] jobs based on the search result of "data scientists" and 5 locations "Singapore", "United States", "United Kingdom", "Canada", "Australia" and "China". The data was collected from 6th June 2022 to 31 June 2022. 

For each job, the following fields are collected:
* jobId
* title
* description
* datePosted
* validThrough: expiring date of listing
* employmentType: type of employment (full time, part time, contract etc)
* hiringOrganization: the hriing company
* hiringOrganizationType: [Can be either hiring company or a recruitment agency]
* industry
* country
* jobLocationType: onsite or remote
* monthsOfExperience: number of months working experience required
* educationRequirements