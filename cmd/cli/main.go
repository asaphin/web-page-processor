package main

import (
	"flag"
	"fmt"
	"github.com/asaphin/web-page-processor/infrastructure/clients"
	"strings"
)

func main() {
	var startPageURL string

	flag.StringVar(&startPageURL, "page", "https://www.rewater.io/", "Web page to start processing")

	flag.Parse()

	getter := clients.NewFastURLGetter()

	h, err := getter.Get(startPageURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nTitle:")
	fmt.Println(h.Title())

	fmt.Println("\nLanguage:")
	fmt.Println(h.Language())

	fmt.Println("\nDescription:")
	fmt.Println(h.Description())

	fmt.Println("\nLinks:")
	links := h.Links()

	for _, link := range links {
		fmt.Println(link)
	}

	fmt.Println("\nMeta:")
	for _, meta := range h.Meta() {
		for key, value := range meta {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println()
	}

	fmt.Println("\nTable of contents:")
	for _, header := range h.TableOfContents() {
		fmt.Printf("%d. %s\n", header.Level, header.Title)
	}

	fmt.Println("\nResponse headers:")
	for name, values := range h.ResponseHeaders() {
		fmt.Printf("%s: %s\n", name, strings.Join(values, ", "))
	}
}
