package shopware

import (
	"context"
	"github.com/go-resty/resty/v2"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type OrderService interface {
	SearchOrdersByTimeRange(ctx context.Context, field string, gte, lte time.Time, page int) ([]Order, error)
	SearchOrdersByIDs(ctx context.Context, IDs []string) ([]Order, error)
	SearchDeliveriesByTimeRange(ctx context.Context, field string, gte, lte time.Time, page int) ([]OrderDelivery, error)
	SearchTransactionsByTimeRange(ctx context.Context, field string, gte, lte time.Time, page int) ([]OrderTransaction, error)
}

type orderService struct {
	client        *resty.Client
	tokenProvider TokenProvider
}

func NewOrderService(client *resty.Client, provider TokenProvider) OrderService {
	return &orderService{
		client:        client,
		tokenProvider: provider,
	}
}

func (s *orderService) SearchOrdersByTimeRange(ctx context.Context, field string, gte, lte time.Time, page int) ([]Order, error) {
	path := "/api/v3/search/order"

	type request struct {
		Page         int          `json:"page"`
		Limit        int          `json:"limit"`
		Filters      []filter     `json:"filter"`
		Associations associations `json:"associations"`
	}

	body := request{
		Page:  page,
		Limit: MaxSearchLimit,
		Filters: []filter{{
			Type:  filterTypeRange,
			Field: field,
			Parameters: &filterParameters{
				GTE: gte.Format(timeFormat),
				LTE: lte.Format(timeFormat),
			},
		}},
		Associations: orderAssociations(),
	}

	var result struct {
		Total int     `json:"total"`
		Data  []Order `json:"data"`
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeaders(s.headers()).
		SetBody(body).
		SetResult(&result).
		Post(path)

	if err := checkHttpResp(resp, err); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (s *orderService) SearchOrdersByIDs(ctx context.Context, IDs []string) ([]Order, error) {
	path := "/api/v3/search/order"

	type request struct {
		Page         int          `json:"page"`
		Limit        int          `json:"limit"`
		Filters      []filter     `json:"filter"`
		Associations associations `json:"associations"`
	}

	body := request{
		Page:  1,
		Limit: MaxSearchLimit,
		Filters: []filter{{
			Type:  filterTypeEqualsAny,
			Field: "id",
			Value: IDs,
		}},
		Associations: orderAssociations(),
	}

	var result struct {
		Total int     `json:"total"`
		Data  []Order `json:"data"`
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeaders(s.headers()).
		SetBody(body).
		SetResult(&result).
		Post(path)

	if err := checkHttpResp(resp, err); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (s *orderService) SearchDeliveriesByTimeRange(ctx context.Context, field string, gte, lte time.Time, page int) ([]OrderDelivery, error) {
	path := "/api/v3/search/order-delivery"

	type request struct {
		Page    int      `json:"page"`
		Limit   int      `json:"limit"`
		Filters []filter `json:"filter"`
	}

	body := request{
		Page:  page,
		Limit: MaxSearchLimit,
		Filters: []filter{{
			Type:  filterTypeRange,
			Field: field,
			Parameters: &filterParameters{
				GTE: gte.Format(timeFormat),
				LTE: lte.Format(timeFormat),
			},
		}},
	}

	var result struct {
		Total int             `json:"total"`
		Data  []OrderDelivery `json:"data"`
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeaders(s.headers()).
		SetBody(body).
		SetResult(&result).
		Post(path)

	if err := checkHttpResp(resp, err); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (s *orderService) SearchTransactionsByTimeRange(ctx context.Context, field string, gte, lte time.Time, page int) ([]OrderTransaction, error) {
	path := "/api/v3/search/order-transaction"

	type request struct {
		Page    int      `json:"page"`
		Limit   int      `json:"limit"`
		Filters []filter `json:"filter"`
	}

	body := request{
		Page:  page,
		Limit: MaxSearchLimit,
		Filters: []filter{{
			Type:  filterTypeRange,
			Field: field,
			Parameters: &filterParameters{
				GTE: gte.Format(timeFormat),
				LTE: lte.Format(timeFormat),
			},
		}},
	}

	var result struct {
		Total int                `json:"total"`
		Data  []OrderTransaction `json:"data"`
	}

	resp, err := s.client.R().
		SetContext(ctx).
		SetHeaders(s.headers()).
		SetBody(body).
		SetResult(&result).
		Post(path)

	if err := checkHttpResp(resp, err); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (s *orderService) headers() map[string]string {
	return tokenHeaders(s.tokenProvider.GetToken())
}

type associations struct {
	Deliveries   []interface{} `json:"deliveries"`
	Transactions []interface{} `json:"transactions"`
	Documents    []interface{} `json:"documents"`
	LineItems    []interface{} `json:"lineItems"`
}

func orderAssociations() associations {
	return associations{
		Deliveries:   []interface{}{},
		Transactions: []interface{}{},
		Documents:    []interface{}{},
		LineItems:    []interface{}{},
	}
}
