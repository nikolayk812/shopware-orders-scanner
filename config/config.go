package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
	"os"
)

// Parse optionally reads from the file if it exists and sets environment variables if they're unset
// Then reading from the env vars parses them into the Config.
func Parse(configFile string, v interface{}) error {
	switch info, err := os.Stat(configFile); {
	case os.IsNotExist(err):
		zap.S().Warnf("file [%s] doesn't exist", configFile)
	case info.IsDir():
		zap.S().Warnf("file [%s] is directory", configFile)
	default:
		if err := gotenv.Load(configFile); err != nil {
			return fmt.Errorf("gotenv.Load [%s]: %w", configFile, err)
		}
	}

	err := envconfig.Process("", v)
	if err != nil {
		return fmt.Errorf("envconfig.Process: %w", err)
	}
	return nil

}

type Shopware struct {
	BaseURL      string `envconfig:"SHOPWARE_BASE_URL" required:"true"`
	ClientID     string `envconfig:"SHOPWARE_CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"SHOPWARE_CLIENT_SECRET" required:"true"`
}

type SendGrid struct {
	Enabled   bool   `envconfig:"SENDGRID_ENABLED" default:"false"`
	APIKey    string `envconfig:"SENDGRID_API_KEY" required:"false"`
	ToEmail   string `envconfig:"SENDGRID_TO_EMAIL" required:"false"`
	ToName    string `envconfig:"SENDGRID_TO_NAME" required:"false"`
	Subject   string `envconfig:"SENDGRID_SUBJECT" required:"false"`
	FromEmail string `envconfig:"SENDGRID_FROM_EMAIL" required:"false"`
	FromName  string `envconfig:"SENDGRID_FROM_NAME" required:"false"`
}
