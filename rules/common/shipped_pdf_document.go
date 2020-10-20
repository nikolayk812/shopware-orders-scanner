package common

import (
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/domains"
)

type ShippedPdfDocument struct{}

func (_ ShippedPdfDocument) Apply(order shopware.Order) (bool, error) {
	// pre-condition
	d, ok := domains.FirstDelivery(order)
	if !ok {
		return false, nil
	}
	if d.StateMachineState.Name != shopware.OrderDeliveryStateShipped {
		return false, nil
	}

	// check
	doc, ok := domains.FirstDocument(order)
	if !ok {
		return false, fmt.Errorf("no document")
	}

	if doc.FileType != "pdf" {
		return false, fmt.Errorf("wrong document file type [%s]", doc.FileType)
	}

	return true, nil
}
