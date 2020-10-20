package render

import (
	"bytes"
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/domains"
	"html/template"
)

type Renderer struct {
	templatePath    string
	shopwareBaseURL string
}

func NewRenderer(templatePath, shopwareBaseURL string) Renderer {
	return Renderer{
		templatePath:    templatePath,
		shopwareBaseURL: shopwareBaseURL,
	}
}

func (r Renderer) RenderHTML(orders []domains.OrderResult, scanned int) ([]byte, error) {
	params := struct {
		BaseURL  string
		Orders   []domains.OrderResult
		Scanned  int
		Detected int
	}{
		BaseURL:  r.shopwareBaseURL,
		Orders:   orders,
		Scanned:  scanned,
		Detected: len(orders),
	}

	t, err := template.ParseFiles(r.templatePath)
	if err != nil {
		return nil, fmt.Errorf("template.ParseFiles [%s] : %w", r.templatePath, err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, params); err != nil {
		return nil, fmt.Errorf("t.Execute : %w", err)
	}

	return body.Bytes(), nil
}
