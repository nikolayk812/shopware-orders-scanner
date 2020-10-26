package consumers

import (
	"github.com/nikolayk812/shopware-orders-scanner/domain"
)

type Consumer interface {
	Consume(orders []domain.OrderResult, scanned int) (Result, error)
}

type Result struct {
	Bytes []byte
}
