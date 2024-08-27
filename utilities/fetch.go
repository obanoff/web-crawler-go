package utilities

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

type Config struct {
	Pages     map[string]uint
	MaxPages  uint
	BaseURL   *url.URL
	mu        *sync.Mutex
	semaphore chan struct{}
	Wg        *sync.WaitGroup
}

func NewConfig(baseURL *url.URL, pool uint8, maxPages uint) *Config {
	if pool == 0 {
		pool = 1
	}
	return &Config{
		Pages:     make(map[string]uint),
		MaxPages:  maxPages,
		BaseURL:   baseURL,
		mu:        &sync.Mutex{},
		semaphore: make(chan struct{}, pool),
		Wg:        &sync.WaitGroup{},
	}
}

func (cfg *Config) PrintReport() {
	var slice []struct {
		page  string
		count uint
	}

	for k, v := range cfg.Pages {
		slice = append(slice, struct {
			page  string
			count uint
		}{
			page:  k,
			count: v,
		})
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].count > slice[j].count
	})

	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, cfg.BaseURL.String())

	for _, p := range slice {
		fmt.Printf("Found %v internal links to %s\n", p.count, p.page)
	}
}

func (cfg *Config) CrawlPage(rawCurrentURL string) error {
	cfg.mu.Lock()
	if cfg.MaxPages <= uint(len(cfg.Pages)) {
		cfg.mu.Unlock()
		return nil
	}
	cfg.mu.Unlock()

	current, err := url.Parse(rawCurrentURL)
	if err != nil {
		return err
	}

	if !current.IsAbs() {
		current = cfg.BaseURL.ResolveReference(current)
	}

	if cfg.BaseURL.Hostname() != current.Hostname() {
		return nil
	}

	normalized, err := NormalizeURL(current.String())
	if err != nil {
		return err
	}

	cfg.mu.Lock()
	if _, ok := cfg.Pages[normalized]; ok {
		cfg.Pages[normalized]++
		cfg.mu.Unlock()
		return nil
	}

	cfg.Pages[normalized] = 1
	cfg.mu.Unlock()

	html, err := GetHTML(current.String())
	if err != nil {
		return err
	}

	urls, err := GetURLsFromHTML(html, cfg.BaseURL.String())
	if err != nil {
		return err
	}

	for _, url := range urls {
		cfg.Wg.Add(1)
		go func(url string) {
			defer func() {
				<-cfg.semaphore
				cfg.Wg.Done()
			}()

			cfg.semaphore <- struct{}{}

			err := cfg.CrawlPage(url)
			if err != nil {
				fmt.Printf("\033[33m%s: %s\033[0m\n", url, err)
			}
		}(url)
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
