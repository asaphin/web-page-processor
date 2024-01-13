package clients

import (
	"fmt"
	"github.com/asaphin/web-page-processor/app"
	"github.com/asaphin/web-page-processor/models"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type FastURLGetter struct {
}

func (g *FastURLGetter) Get(url string) (app.HTML, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(url)

	if err := fasthttp.DoTimeout(request, response, 5*time.Second); err != nil {
		return nil, err
	}

	contentType := response.Header.Peek("Content-Type")
	if !strings.HasPrefix(string(contentType), "text/html") {
		return nil, fmt.Errorf("expected text/html content, got %s", contentType)
	}

	headers := make(map[string][]string)

	response.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = append(headers[string(key)], strings.Split(string(value), ", ")...)
	})

	return models.WrapHTML(url, response.Body(), headers)
}

func NewFastURLGetter() app.URLGetter {
	return &FastURLGetter{}
}
