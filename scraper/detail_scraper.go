package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ScraperLog struct {
	Records map[string]bool
	fh      *os.File
}

func OpenLog(path string) ScraperLog {
	log := ScraperLog{Records: make(map[string]bool)}
	buff, err := os.ReadFile(path)
	if err == nil {
		for _, l := range strings.Split(string(buff), "\n") {
			log.Records[strings.Split(l, ",")[0]] = true
		}
	}
	log.fh, _ = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	return log
}

func (log *ScraperLog) WriteLog(r *ScrapeResult) {
	line := r.link + ", "
	if len(r.errMessage) > 1 {
		line += "ERROR: " + r.errMessage
	} else {
		line += "SUCCESS"
	}
	fmt.Println(line)
	log.fh.WriteString(line + "\n")
	log.fh.Sync()
}

func (log *ScraperLog) CheckLink(link string) (string, error) {
	start := strings.Index(link, "https")
	if start == -1 {
		return link, errors.New("invalid")
	}
	end := strings.LastIndex(link, "?")
	if end == -1 {
		return link, errors.New("invalid")
	}
	l := link[start:end]
	_, ok := log.Records[l]
	if ok {
		return l, errors.New("duplicate")
	} else {
		return l, nil
	}
}

type ScrapeResult struct {
	link       string
	body       []byte
	errMessage string
}

var scraperLog = OpenLog("./data/detail-scraper.log")
var queue []string
var queueIndex = 0
var completed = 0

func ScrapeFolder() {
	files, _ := os.ReadDir("./data")
	for _, f := range files {
		if strings.Contains(f.Name(), ".txt") {
			ScrapeDetail("./data/" + f.Name())
		}
	}
}

func ScrapeDetail(filename string) {
	b, _ := os.ReadFile(filename)
	queue = strings.Split(string(b), "\n")
	queueIndex = -1
	completed = 0
	ch := make(chan ScrapeResult)

	queueIndex++
	go scrape(queueIndex, ch)

	for r := range ch {
		completed++
		queueIndex++
		if queueIndex < len(queue) {
			go scrape(queueIndex, ch)
		}

		if completed == len(queue) {
			os.Rename(filename, strings.Replace(filename, "/data", "/data/backup", 1))
			close(ch)
		}

		if r.errMessage == "duplicate" || r.errMessage == "invalid" {
			continue
		}
		time.Sleep(RandDuration(1, 3))
		n := r.link[strings.LastIndex(r.link, "-")+1:]
		outputPath := fmt.Sprintf("./data/json/%s.json", n)
		reg := regexp.MustCompile(`(?s)<script type=\"application/ld\+json\">([^<]+)?`)
		regMatch := reg.FindStringSubmatch(string(r.body))
		if len(regMatch) == 2 {
			os.WriteFile(outputPath, []byte(regMatch[1]), 0644)
		} else {
			os.WriteFile(strings.ReplaceAll(outputPath, "json", "html"), r.body, 0644)
			r.errMessage = "JSON Not Found"
		}
		scraperLog.WriteLog(&r)
	}
}

func scrape(index int, ch chan ScrapeResult) {

	if index >= len(queue) {
		return
	}
	link, err := scraperLog.CheckLink(queue[index])
	result := ScrapeResult{link: link}
	if err != nil {
		result.errMessage = err.Error()
		ch <- result
		return
	}

	resp, err := http.Get(link)

	if err != nil {
		result.errMessage = err.Error()
		ch <- result
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		result.errMessage = err.Error()
		ch <- result

	} else {
		result.body = body
		ch <- result
	}
}
