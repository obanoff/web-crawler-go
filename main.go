package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/obanoff/web-crawler-go/utilities"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	// default values
	var (
		pool     uint8 = 10
		maxPages uint  = 50
	)

	if len(os.Args) > 2 {
		if len(os.Args) > 4 {
			fmt.Println("too many arguments provided")
			os.Exit(1)
		}

		v1, err := strconv.Atoi(os.Args[2])
		if err != nil || v1 < 0 {
			fmt.Println("negative thread pool")
			os.Exit(1)
		}

		v2, err := strconv.Atoi(os.Args[3])
		if err != nil || v2 < 0 {
			fmt.Println("negative max pages")
			os.Exit(1)
		}

		pool, maxPages = uint8(v1), uint(v2)
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

	cfg := utilities.NewConfig(url, pool, maxPages)

	err = cfg.CrawlPage(url.String())
	if err != nil {
		fmt.Println(err)
		// os.Exit(1)
	}

	cfg.Wg.Wait()

	cfg.PrintReport()
}
