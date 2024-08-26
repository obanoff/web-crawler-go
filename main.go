package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/obanoff/web-crawler-go/utilities"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	url, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println("error parsing url")
		os.Exit(1)
	}

	if !url.IsAbs() {
		fmt.Println("invalid url provided")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", url.String())

	pages := make(map[string]int)

	err = utilities.CrawlPage(url.String(), url.String(), pages)
	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
	}

	for k, v := range pages {
		fmt.Printf("%s: %v\n", k, v)
	}

}
