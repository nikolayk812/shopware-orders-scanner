package domains

import (
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"sort"
)

func FirstDelivery(order shopware.Order) (shopware.OrderDelivery, bool) {
	if len(order.Deliveries) == 0 {
		return shopware.OrderDelivery{}, false
	}

	return order.Deliveries[0], true
}

func FirstDocument(order shopware.Order) (shopware.OrderDocument, bool) {
	if len(order.Documents) == 0 {
		return shopware.OrderDocument{}, false
	}

	return order.Documents[0], true
}

func LatestTransaction(order shopware.Order) (shopware.OrderTransaction, bool) {
	transactions := order.Transactions
	if len(transactions) == 0 {
		return shopware.OrderTransaction{}, false
	}

	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.Before(transactions[j].CreatedAt)
	})

	//sorted in ASC order, take the latest then
	return transactions[len(transactions)-1], true
}

func TrackingCode(order shopware.Order) string {
	if len(order.Deliveries) == 0 {
		return "absent"
	}

	if len(order.Deliveries[0].TrackingCodes) == 0 {
		return "absent"
	}

	return order.Deliveries[0].TrackingCodes[0]
}
