// Golang tour exercise: Web Crawler
package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type FetcherResult struct{
	body string
	urls []string
	err error
	depth int
}

func _Crawl(url string, depth int, fetcher Fetcher, ch chan *FetcherResult) {
	var fr FetcherResult 
	if depth <= 0 {
		ch <- nil
		return
	}
		
	fr.body, fr.urls, fr.err = fetcher.Fetch(url)
	fr.depth = depth
	if fr.err != nil {
		fmt.Println(fr.err)
		ch <- nil
		return
	} else {
		fmt.Printf("found: %s %q\n", url, fr.body)
	}	
	ch <- &fr 
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
// depth is how far the current url is from the first fetched url
func Crawl(url string, depth int, fetcher Fetcher) {
	ch := make(chan *FetcherResult)
	fetched := make(map[string]bool)
	go _Crawl(url, depth, fetcher, ch)
	fetched[url] = true
	goRtn := 1 // track number of goroutines that are fetching 
	for goRtn > 0 {
		fr := <-ch
		goRtn-- // Get fetched result from a goroutine which is finished
		if fr == nil {
			continue
		}
		for _, childUrl := range fr.urls {
			if _, done := fetched[childUrl]; !done {
				goRtn++
				fetched[childUrl] = true
				go _Crawl(childUrl, fr.depth - 1, fetcher, ch)
			}
		}
	}
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher {
	"http://golang.org/": &fakeResult {
		"The Go Programming Language",
		[]string {
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult {
		"Packages",
		[]string {
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult {
		"Package fmt",
		[]string {
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult {
		"Package os",
		[]string {
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
