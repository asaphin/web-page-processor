package models

import (
	"bytes"
	"github.com/asaphin/web-page-processor/app"
	"github.com/asaphin/web-page-processor/domain"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

type HTMLWrapper struct {
	url      *url.URL
	doc      *html.Node
	rawBytes []byte
	headers  map[string][]string
}

func WrapHTML(rawURL string, data []byte, headers map[string][]string) (app.HTML, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &HTMLWrapper{
		url:      parsedURL,
		doc:      doc,
		rawBytes: data,
		headers:  headers,
	}, nil
}

func (h *HTMLWrapper) Title() string {
	var title string

	dfs(h.doc, func(n *html.Node) (terminate bool) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n.FirstChild.Data
			return true
		}

		return false
	})

	return title
}

func (h *HTMLWrapper) Language() string {
	var language string

	dfs(h.doc, func(n *html.Node) (terminate bool) {
		if n.Type == html.ElementNode && n.Data == "html" {
			for _, attribute := range n.Attr {
				if attribute.Key == "lang" {
					language = attribute.Val
					return true
				}
			}
		}

		return false
	})

	return language
}

func (h *HTMLWrapper) Links() []string {
	var links []string

	dfs(h.doc, func(n *html.Node) (terminate bool) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attribute := range n.Attr {
				if attribute.Key == "href" {
					link := attribute.Val
					if !strings.HasPrefix(link, "http") {
						newURL := &url.URL{
							Scheme: h.url.Scheme,
							Host:   h.url.Host,
						}

						link = newURL.String() + link
					}

					links = append(links, link)
				}
			}
		}

		return false
	})

	return links
}

func (h *HTMLWrapper) Bytes() []byte {
	return h.rawBytes
}

func (h *HTMLWrapper) Meta() []map[string]string {
	var metaParamsList []map[string]string

	dfs(h.doc, func(n *html.Node) (terminate bool) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			metaParams := make(map[string]string)
			for _, attribute := range n.Attr {
				metaParams[attribute.Key] = attribute.Val
			}

			if len(metaParams) > 0 {
				metaParamsList = append(metaParamsList, metaParams)
			}
		}

		return false
	})

	return metaParamsList
}

func (h *HTMLWrapper) TableOfContents() []domain.Header {
	tableOfContents := make([]domain.Header, 0)

	dfs(h.doc, func(n *html.Node) (terminate bool) {
		if n.Type == html.ElementNode && strings.HasPrefix(n.Data, "h") && len(n.Data) == 2 && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			level := int(n.Data[1] - '0') // Extract header level (H1, H2, ..., H6)
			title := strings.TrimSpace(n.FirstChild.Data)

			header := domain.Header{
				Level: level,
				Title: title,
			}

			tableOfContents = append(tableOfContents, header)
		}

		return false
	})

	return tableOfContents
}

func (h *HTMLWrapper) ResponseHeaders() map[string][]string {
	return h.headers
}

// bfs traverses the HTML tree in a breadth-first manner and applies the given function to each node.
func bfs(root *html.Node, visit func(*html.Node) bool) {
	if root == nil {
		return
	}

	queue := []*html.Node{root}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		terminate := visit(currentNode)
		if terminate {
			return
		}

		for child := currentNode.FirstChild; child != nil; child = child.NextSibling {
			queue = append(queue, child)
		}
	}
}

// dfs traverses the HTML tree in a depth-first manner and applies the given function to each node.
func dfs(root *html.Node, visit func(*html.Node) bool) {
	if root == nil {
		return
	}

	var recursiveDFS func(*html.Node)
	recursiveDFS = func(node *html.Node) {
		terminate := visit(node)
		if terminate {
			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			recursiveDFS(child)
		}
	}

	recursiveDFS(root)
}
