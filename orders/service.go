package orders

import (
	"context"
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/checks"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/domain"
	"go.uber.org/zap"
	"sort"
	"time"
)

type Service struct {
	orderCli shopware.OrderService
	engine   checks.Engine
}

func NewService(orderCli shopware.OrderService, engine checks.Engine) Service {
	return Service{
		orderCli: orderCli,
		engine:   engine,
	}
}

type FilterRequest struct {
	From, To                  time.Time
	IncludeCreated            bool
	IncludeUpdated            bool
	IncludeDeliveryUpdated    bool
	IncludeTransactionUpdated bool
}

// returns only bad orders
func (s Service) ScanOrders(ctx context.Context, req FilterRequest) ([]domain.OrderResult, int, error) {
	var result []domain.OrderResult
	processed := map[string]bool{}

	if req.IncludeUpdated {
		orders, err := s.getAllOrders(ctx, s.orderCli.SearchByTimeRange, "updatedAt", req.From, req.To)
		if err != nil {
			return nil, 0, fmt.Errorf("SearchByTimeRange(updatedAt) : %w", err)
		}
		badOrders := s.processOrders(orders, processed)
		result = append(result, badOrders...)
	}

	if req.IncludeCreated {
		orders, err := s.getAllOrders(ctx, s.orderCli.SearchByTimeRange, "createdAt", req.From, req.To)
		if err != nil {
			return nil, 0, fmt.Errorf("SearchByTimeRange(createdAt) : %w", err)
		}
		badOrders := s.processOrders(orders, processed)
		result = append(result, badOrders...)
	}

	if req.IncludeDeliveryUpdated {
		orders, err := s.getAllOrders(ctx, s.searchOrdersByDeliveries, "updatedAt", req.From, req.To)
		if err != nil {
			return nil, 0, fmt.Errorf("searchOrdersByDeliveries : %w", err)
		}
		badOrders := s.processOrders(orders, processed)
		result = append(result, badOrders...)
	}

	if req.IncludeTransactionUpdated {
		orders, err := s.getAllOrders(ctx, s.searchOrdersByTransactions, "updatedAt", req.From, req.To)
		if err != nil {
			return nil, 0, fmt.Errorf("searchOrdersByTransactions : %w", err)
		}
		badOrders := s.processOrders(orders, processed)
		result = append(result, badOrders...)
	}

	sortResult(result)
	return result, len(processed), nil
}

func (s Service) getAllOrders(ctx context.Context,
	searchFunc func(context.Context, string, time.Time, time.Time, int) ([]shopware.Order, error),
	field string, from, to time.Time) ([]shopware.Order, error) {

	var result []shopware.Order
	page := 1
	for {
		pageOrders, err := searchFunc(ctx, field, from, to, page)
		if err != nil {
			return nil, fmt.Errorf("failed to get orders : %w", err)
		}
		result = append(result, pageOrders...)
		page++

		if len(pageOrders) < shopware.MaxSearchLimit {
			break
		}
	}
	zap.S().Infof("got %d orders by %s", len(result), field)
	return result, nil
}

func (s Service) searchOrdersByDeliveries(ctx context.Context, _ string, from, to time.Time, page int) ([]shopware.Order, error) {
	deliveries, err := s.orderCli.SearchDeliveriesByTimeRange(ctx, "updatedAt", from, to, page)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries by updatedAt : %w", err)
	}

	if len(deliveries) == 0 {
		return []shopware.Order{}, nil
	}

	var orderIDs []string
	for _, d := range deliveries {
		orderIDs = append(orderIDs, d.OrderID)
	}
	return s.orderCli.SearchByIDs(ctx, orderIDs)
}

func (s Service) searchOrdersByTransactions(ctx context.Context, _ string, from, to time.Time, page int) ([]shopware.Order, error) {
	txs, err := s.orderCli.SearchTransactionsByTimeRange(ctx, "updatedAt", from, to, page)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by updatedAt : %w", err)
	}

	if len(txs) == 0 {
		return []shopware.Order{}, nil
	}

	var orderIDs []string
	for _, d := range txs {
		orderIDs = append(orderIDs, d.OrderID)
	}
	return s.orderCli.SearchByIDs(ctx, orderIDs)
}

// returns only bad orders
func (s Service) processOrders(orders []shopware.Order, processed map[string]bool) []domain.OrderResult {
	var badOrders []domain.OrderResult
	for _, order := range orders {
		if processed[order.ID] {
			continue
		}
		processed[order.ID] = true

		errors := s.engine.ProcessOrder(order)
		if len(errors) > 0 {
			badOrders = append(badOrders, domain.OrderResult{
				OrderID:      order.ID,
				OrderNumber:  order.Number,
				ChannelID:    order.SalesChannelID,
				TrackingCode: domain.TrackingCode(order),
				CreatedDate:  trySubString(order.CreatedAt, 10),
				Errors:       errors,
			})
		}
	}
	return badOrders
}

func sortResult(r []domain.OrderResult) {
	sort.SliceStable(r, func(i, j int) bool {
		if len(r[i].Errors) != len(r[j].Errors) {
			return len(r[i].Errors) < len(r[j].Errors)
		}

		for iKey := range r[i].Errors {
			_, ok := r[j].Errors[iKey]
			if !ok {
				for jKey := range r[j].Errors {
					return iKey < jKey
				}
			}
		}

		if r[i].CreatedDate != r[j].CreatedDate {
			return r[i].CreatedDate < r[j].CreatedDate
		}

		return r[i].OrderNumber < r[j].OrderNumber
	})
}

func trySubString(s string, l int) string {
	if len(s) <= l {
		return s
	}
	return s[:l]
}
