package main

import (
	"flag"
	"fmt"
	"github.com/asaphin/web-page-processor/infrastructure/clients"
	"strings"
)

func main() {
	var startPageURL string

	flag.StringVar(&startPageURL, "page", "https://www.youtube.com/", "Web page to start processing")

	flag.Parse()

	getter := clients.NewFastURLGetter()

	h, err := getter.Get(startPageURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(h.Title())
	fmt.Println(h.Language())

	links := h.Links()

	for _, link := range links {
		fmt.Println(link)
	}

	for _, meta := range h.Meta() {
		for key, value := range meta {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println()
	}

	for _, header := range h.TableOfContents() {
		fmt.Printf("%d. %s\n", header.Level, header.Title)
	}

	for name, values := range h.ResponseHeaders() {
		fmt.Printf("%s: %s\n", name, strings.Join(values, ", "))
	}
}
