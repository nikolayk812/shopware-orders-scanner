package shopware

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type ProductService interface {
	SearchProductByNumber(ctx context.Context, prodNumber string) (Product, error)
}

type productService struct {
	client        *resty.Client
	tokenProvider TokenProvider
}

func NewProductService(client *resty.Client, provider TokenProvider) ProductService {
	return &productService{
		client:        client,
		tokenProvider: provider,
	}
}

func (s *productService) SearchProductByNumber(ctx context.Context, prodNumber string) (Product, error) {
	path := "/api/v1/product?filter[product.productNumber]=" + prodNumber

	var result struct {
		Total int       `json:"total"`
		Data  []Product `json:"data"`
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeaders(s.headers()).
		SetResult(&result).
		Get(path)

	if err := checkHttpResp(resp, err); err != nil {
		return Product{}, err
	}

	if result.Total != 1 {
		return Product{}, fmt.Errorf("unexpected number of products: %d", result.Total)
	}

	return result.Data[0], nil
}

func (s *productService) headers() map[string]string {
	return tokenHeaders(s.tokenProvider.GetToken())
}
