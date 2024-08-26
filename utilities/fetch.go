package utilities

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) error {
	base, _ := url.Parse(rawBaseURL)
	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		return err
	}

	if !current.IsAbs() {
		current = base.ResolveReference(current)
	}

	if base.Hostname() != current.Hostname() {
		return nil
	}

	normalized, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		return err
	}

	if _, ok := pages[normalized]; ok {
		pages[normalized]++
		return nil
	}

	pages[normalized] = 1

	fmt.Printf("crawling %s\n", current.String())

	html, err := GetHTML(current.String())
	if err != nil {
		return err
	}

	urls, err := GetURLsFromHTML(html, rawBaseURL)
	if err != nil {
		return err
	}

	for _, url := range urls {
		err = CrawlPage(rawBaseURL, url, pages)
		if err != nil {
			log.Printf("Failed to crawl %s: %v\n", url, err)
			continue
		}
	}
	return nil
}

func GetHTML(rawURL string) (string, error) {
	resp, err := http.DefaultClient.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("%s: %v", resp.Status, resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", fmt.Errorf("not html response")
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}
