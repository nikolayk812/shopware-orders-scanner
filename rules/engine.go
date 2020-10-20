package rules

import (
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"go.uber.org/zap"
)

type Engine struct {
	rules map[string]Rule
}

func NewEngine(rules map[string]Rule) Engine {
	return Engine{rules: rules}
}

func (e Engine) ProcessOrder(order shopware.Order) map[string]error {
	errors := map[string]error{}
	for ruleName, rule := range e.rules {
		_, err := rule.Apply(order)
		if err != nil {
			errors[ruleName] = err
		}
		processResult(err, order, ruleName)
	}
	return errors
}

func processResult(err error, order shopware.Order, ruleName string) {
	if err != nil {
		zap.S().Errorf("order [%s] has failed to pass [%s] check : %v", order.ID, ruleName, err)
		return
	}

	zap.S().Debugf("order [%s] has passed [%s] check", order.ID, ruleName)
}
