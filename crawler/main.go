package main

import (
	"fmt"
	"os"
	"time"
)

const usage = `
usage:
	crawler <starting-url>
`

func worker(linkq chan string, resultsq chan []string) {
	for link := range linkq {
		plinks, err := getPageLinks(link)
		if err != nil {
			fmt.Printf("ERROR fetching %s: %v\n", link, err)
			continue
		}
		fmt.Printf("%s (%d links)\n", link, len(plinks.Links)) // len(...) = # of links found on the page
		time.Sleep(time.Millisecond * 500)
		if len(plinks.Links) > 0 {
			go func(links []string) { // run separately from the worker
				resultsq <- plinks.Links
			}(plinks.Links)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	nWorders := 50                   // faster w more workers as well, since working concurrently
	linkq := make(chan string, 1000) // bigger the number, the faster bc doesn't get blocked as much
	resultsq := make(chan []string, 1000)
	for i := 0; i < nWorders; i++ {
		go worker(linkq, resultsq)
	}
	linkq <- os.Args[1]

	seen := map[string]bool{}
	for links := range resultsq {
		for _, link := range links { // links is []string
			if !seen[link] {
				seen[link] = true
				linkq <- link
			}
		}
	}
}
