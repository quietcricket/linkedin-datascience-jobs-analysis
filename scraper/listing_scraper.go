package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gosimple/slug"
)

var scripts map[string]string = loadJs()

func main() {
	// Remove headless mode, easier to observe
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.WindowSize(1280, 800), chromedp.Flag("headless", false))

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel1()

	browserCtx, cancel2 := chromedp.NewContext(allocCtx)
	defer cancel2()

	search("Data Scientist", "United States", browserCtx)
	search("Data Scientist", "United Kingdom", browserCtx)
	search("Data Scientist", "Canada", browserCtx)
	search("Data Scientist", "Singapore", browserCtx)
	search("Data Scientist", "Australia", browserCtx)
	search("Data Scientist", "China", browserCtx)
}

func search(keyword, location string, ctx context.Context) {

	// open a file to store the data
	t := time.Now()
	filename := fmt.Sprintf("./data/%d%02d%02d-%s-%s.txt", t.Year(), t.Month(), t.Day(), slug.Make(keyword), slug.Make(location))
	output, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 644)
	defer output.Close()

	// build the search url
	params := url.Values{}
	params.Add("keywords", keyword)
	params.Add("location", location)
	url := "https://www.linkedin.com/jobs/search?" + params.Encode()

	// open browser and wait for the page is loaded
	chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(".jobs-search__results-list"))

	result := make(chan string)
	defer close(result)
	for {
		go scrapeLinks(ctx, result)
		line := <-result
		if line == "stop" {
			break
		} else if len(line) > 0 {
			output.WriteString(line)
		}
	}
}

func loadJs() map[string]string {
	m := make(map[string]string)

	for _, f := range []string{"scrape", "scroll-listing", "show-more"} {
		content, _ := os.ReadFile(fmt.Sprintf("./js/%s.js", f))
		m[f] = string(content)
	}
	return m
}

func randDelay(low, high float64) time.Duration {
	return time.Duration((rand.Float64()*(high-low) + low) * float64(time.Second))
}

func scrapeLinks(ctx context.Context, ch chan string) {
	var res []byte
	chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools(scripts["scrape"], &res, chromedp.EvalAsValue),
		chromedp.EvaluateAsDevTools(scripts["scroll-listing"], nil),
	)

	time.Sleep(randDelay(0.7, 1.5))

	chromedp.Run(ctx, chromedp.EvaluateAsDevTools(scripts["show-more"], nil))

	time.Sleep(randDelay(0.7, 1.5))

	result, _ := strconv.Unquote(string(res))
	fmt.Printf("Links: %s", result)
	ch <- result
}
