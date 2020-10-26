package mail

import (
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/config"
	"github.com/nikolayk812/shopware-orders-scanner/consumers"
	"github.com/nikolayk812/shopware-orders-scanner/consumers/html"
	"github.com/nikolayk812/shopware-orders-scanner/domain"
	ms "github.com/nikolayk812/shopware-orders-scanner/mail"
)

type Sender struct {
	sender   ms.Sender
	renderer html.Renderer
}

func NewSender(swConf config.Shopware, sgConf config.SendGrid) Sender {
	return Sender{
		sender:   ms.NewSender(sgConf),
		renderer: html.NewRenderer("./consumers/html/template.twig", swConf.BaseURL),
	}
}

func (s Sender) Consume(orders []domain.OrderResult, scanned int) (consumers.Result, error) {
	out, err := s.renderer.Consume(orders, scanned)
	if err != nil {
		return consumers.Result{}, fmt.Errorf("renderer.Consume : %w", err)
	}

	if err := s.sender.SendMail(string(out.Bytes)); err != nil {
		return consumers.Result{}, fmt.Errorf("sender.SendMail : %w", err)
	}

	return consumers.Result{}, nil
}
