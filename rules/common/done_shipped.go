package common

import (
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/domains"
)

type DoneDeliveryNotOpen struct{}

func (_ DoneDeliveryNotOpen) Apply(order shopware.Order) (bool, error) {
	// pre-condition
	if order.StateMachineState.Name != shopware.OrderStateDone {
		return false, nil
	}

	// check
	d, ok := domains.FirstDelivery(order)
	if !ok {
		return false, fmt.Errorf("no deliveries")
	}
	if d.StateMachineState.Name == shopware.OrderDeliveryStateOpen {
		return false, fmt.Errorf("wrong delivery state [%s] when shipped",
			d.StateMachineState.Name)
	}

	return true, nil
}
