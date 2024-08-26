package utilities

import (
	"errors"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func NormalizeURL(rawURL string) (string, error) {
	re, err := regexp.Compile(`^[A-z]+://|/$`)
	if err != nil {
		return "", errors.New("error compiling regular expression")
	}

	return re.ReplaceAllString(rawURL, ""), nil
}

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	parsedURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(htmlBody)

	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0)

	err = traverse(doc, parsedURL, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func traverse(node *html.Node, base *url.URL, urls *[]string) error {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				newURL, err := url.Parse(a.Val)
				if err != nil {
					return err
				}

				if !newURL.IsAbs() {
					newURL = base.ResolveReference(newURL)
				}

				if str := newURL.String(); !slices.Contains(*urls, str) {
					*urls = append(*urls, str)
				}
				break
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := traverse(c, base, urls); err != nil {
			return err
		}
	}

	return nil
}
