package main

import (
	"fmt"
	"sync"
)

type mutexCache struct {
	mux   sync.Mutex
	store map[string]bool
}

func (cache *mutexCache) setVisited(name string) bool {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	if cache.store[name] {
		return true
	}

	cache.store[name] = true

	return false
}

var cacheInstance = mutexCache{store: make(map[string]bool)}

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func crawlInner(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	if cacheInstance.setVisited(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	for _, u := range urls {
		wg.Add(1)

		go crawlInner(u, depth-1, fetcher, wg)
	}

	return
}

func Crawl(url string, depth int, fetcher Fetcher) {
	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)

	go crawlInner(url, depth, fetcher, waitGroup)

	waitGroup.Wait()
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

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

var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
