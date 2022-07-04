package main

import "os"

func main() {
	switch os.Args[1] {
	case "listing":
		ScrapeListing()
	case "detail":
		ScrapeFolder()
	case "json":
		ProcessJson()
	}
}
