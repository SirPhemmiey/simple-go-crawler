package crawler

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func SingleFetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func MultipleFetch(urls []string) ([]string, error) {
	var fetchedUrls []string

	for _, u := range urls {
		singleUrl, err := SingleFetch(u)
		if err != nil {
			return fetchedUrls, err
		}
		fetchedUrls = append(fetchedUrls, singleUrl)
	}
	return fetchedUrls, nil
}

// extract links from html content
func ExtractLinks(body string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func MultipleExtractLinks(body []string) ([][]string, error) {
	var links [][]string
	for _, l := range body {
		extractedLinks, err := ExtractLinks(l)
		if err != nil {
			return links, err
		}
		links = append(links, extractedLinks)
	}

	return links, nil
}
