package common

import (
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/domains"
)

type ReturnedRefundedState struct{}

func (_ ReturnedRefundedState) Apply(order shopware.Order) (bool, error) {
	// pre-conditions
	d, ok := domains.FirstDelivery(order)
	if !ok {
		return false, nil
	}

	// check
	tx, ok := domains.LatestTransaction(order)
	if !ok {
		return false, fmt.Errorf("no transactions")
	}

	txState := tx.StateMachineState.Name
	switch d.StateMachineState.Name {
	case shopware.OrderDeliveryStateReturned:
		if txState != shopware.OrderTransactionStateRefunded {
			return false, fmt.Errorf("wrong payment state [%s] expected [%s]",
				txState, shopware.OrderTransactionStateRefunded)
		}
	case shopware.OrderDeliveryStateReturnedPartially:
		if txState != shopware.OrderTransactionStateRefundedPartially {
			return false, fmt.Errorf("wrong payment state [%s] expected [%s]",
				txState, shopware.OrderTransactionStateRefundedPartially)
		}
	default:
		return false, nil //pre-conditions are not met
	}

	return true, nil
}
