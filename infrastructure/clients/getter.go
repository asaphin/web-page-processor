package clients

import (
	"fmt"
	"github.com/asaphin/web-page-processor/app"
	"github.com/asaphin/web-page-processor/models"
	"io"
	"net/http"
	"strings"
)

type StandardURLGetter struct {
}

func (g *StandardURLGetter) Get(url string) (app.HTML, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("expected text/html content, got %s", contentType)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return models.WrapHTML(url, data, response.Header)
}

func NewStandardURLGetter() app.URLGetter {
	return &StandardURLGetter{}
}
