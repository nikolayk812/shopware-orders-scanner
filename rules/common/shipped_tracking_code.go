package common

import (
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/domains"
)

type ShippedTrackingCode struct{}

func (_ ShippedTrackingCode) Apply(order shopware.Order) (bool, error) {
	// pre-condition
	d, ok := domains.FirstDelivery(order)
	if !ok {
		return false, nil
	}
	if d.StateMachineState.Name != shopware.OrderDeliveryStateShipped {
		return false, nil
	}

	// check
	if len(d.TrackingCodes) == 0 {
		return false, fmt.Errorf("no tracking code")
	}

	tc := d.TrackingCodes[0]
	if tc == "" {
		return false, fmt.Errorf("tracking code is empty")
	}

	return true, nil
}
