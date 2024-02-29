# Scraper
This is a go lambda that goes through the proxy api to receive website html. Once received it parses the html and does a keyword check on the job description. If any keyword exists in the description then the job link and company are sent to discord for manual review.

## Prerequisites
You must have the following installed/configured on your system for this to work correctly<br />
1. [Go](https://go.dev/doc/install)
