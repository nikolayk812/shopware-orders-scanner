package orders

import (
	"context"
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/domains"
	"github.com/nikolayk812/shopware-orders-scanner/rules"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	orderCli shopware.OrderService
	engine   rules.Engine
}

func NewService(orderCli shopware.OrderService, engine rules.Engine) Service {
	return Service{
		orderCli: orderCli,
		engine:   engine,
	}
}

// returns only bad orders
func (s Service) ScanOrders(ctx context.Context, from, to time.Time) ([]domains.OrderResult, int, error) {
	processed := map[string]bool{}
	var result []domains.OrderResult

	var orders []shopware.Order
	page := 1
	for {
		pageOrders, err := s.orderCli.SearchOrdersByTimeRange(ctx, "updatedAt", from, to, page)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get orders by updatedAt : %w", err)
		}
		orders = append(orders, pageOrders...)
		page++

		if len(pageOrders) < shopware.MaxSearchLimit {
			break
		}
	}
	zap.S().Infof("got %d orders by updatedAt", len(orders))
	badOrders := s.processOrders(orders, processed)
	result = append(result, badOrders...)

	page = 1
	orders = nil
	for {
		pageOrders, err := s.orderCli.SearchOrdersByTimeRange(ctx, "createdAt", from, to, page)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get orders by createdAt : %w", err)
		}
		orders = append(orders, pageOrders...)
		page++

		if len(pageOrders) < shopware.MaxSearchLimit {
			break
		}
	}
	zap.S().Infof("got %d orders by createdAt", len(orders))
	badOrders = s.processOrders(orders, processed)
	result = append(result, badOrders...)

	page = 1
	orders = nil
	for {
		pageOrders, err := s.searchOrdersByDeliveries(ctx, from, to, page)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get orders by deliveries : %w", err)
		}
		orders = append(orders, pageOrders...)
		page++

		if len(pageOrders) < shopware.MaxSearchLimit {
			break
		}
	}
	zap.S().Infof("got %d orders by deliveries", len(orders))
	badOrders = s.processOrders(orders, processed)
	result = append(result, badOrders...)

	page = 1
	orders = nil
	for {
		pageOrders, err := s.searchOrdersByTransactions(ctx, from, to, page)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get orders by transactions : %w", err)
		}
		orders = append(orders, pageOrders...)
		page++

		if len(pageOrders) < shopware.MaxSearchLimit {
			break
		}
	}
	zap.S().Infof("got %d orders by transactions", len(orders))
	badOrders = s.processOrders(orders, processed)
	result = append(result, badOrders...)

	return result, len(processed), nil
}

func (s Service) searchOrdersByDeliveries(ctx context.Context, from, to time.Time, page int) ([]shopware.Order, error) {
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
	return s.orderCli.SearchOrdersByIDs(ctx, orderIDs)
}

func (s Service) searchOrdersByTransactions(ctx context.Context, from, to time.Time, page int) ([]shopware.Order, error) {
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
	return s.orderCli.SearchOrdersByIDs(ctx, orderIDs)
}

// returns only bad orders
func (s Service) processOrders(orders []shopware.Order, processed map[string]bool) []domains.OrderResult {
	var badOrders []domains.OrderResult
	for _, order := range orders {
		if processed[order.ID] {
			continue
		}
		processed[order.ID] = true

		errors := s.engine.ProcessOrder(order)
		if len(errors) > 0 {
			badOrders = append(badOrders, domains.OrderResult{
				OrderID:      order.ID,
				OrderNumber:  order.Number,
				ChannelID:    order.SalesChannelID,
				TrackingCode: domains.TrackingCode(order),
				Errors:       errors,
			})
		}
	}
	return badOrders
}
