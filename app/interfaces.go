package app

import "github.com/asaphin/web-page-processor/domain"

type URLGetter interface {
	Get(url string) (HTML, error)
}

type HTML interface {
	Title() string
	Description() string
	Language() string
	Links() []string
	Bytes() []byte
	Meta() []map[string]string
	TableOfContents() []domain.Header
	ResponseHeaders() map[string][]string
}
