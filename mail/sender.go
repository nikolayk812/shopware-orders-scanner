package mail

import (
	"fmt"
	"github.com/nikolayk812/shopware-orders-scanner/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sender struct {
	conf config.SendGrid
}

func NewSender(conf config.SendGrid) Sender {
	return Sender{
		conf: conf,
	}
}

func (s Sender) SendMail(htmlContent string) error {
	if !s.conf.Enabled {
		return nil
	}

	from := mail.NewEmail(s.conf.FromName, s.conf.FromEmail)
	to := mail.NewEmail(s.conf.ToName, s.conf.ToEmail)
	message := mail.NewSingleEmail(from, s.conf.Subject, to, "Shopware Orders Scanner", htmlContent)
	client := sendgrid.NewSendClient(s.conf.APIKey)

	resp, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("client.Send: %w", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("failed to send, status : %d, body : [%s]", resp.StatusCode, resp.Body)
	}

	return nil
}
