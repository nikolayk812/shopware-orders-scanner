package checks

import "github.com/nikolayk812/shopware-orders-scanner/clients/shopware"

type Check interface {
	//(true, nil) -> okay
	//(false, nil) -> skipped, i.e. not applicable
	//(false, err) -> failure
	Apply(order shopware.Order) (bool, error)
}
