package main

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/config"
	"github.com/nikolayk812/shopware-orders-scanner/mail"
	"github.com/nikolayk812/shopware-orders-scanner/orders"
	"github.com/nikolayk812/shopware-orders-scanner/render"
	"github.com/nikolayk812/shopware-orders-scanner/rules"
	"github.com/nikolayk812/shopware-orders-scanner/rules/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

type mainConfig struct {
	config.Shopware
	config.SendGrid
}

func main() {
	logger, err := buildLogger()
	if err != nil {
		log.Fatalf("buildLogger: %v", err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	now := time.Now()
	zap.S().Info("starting Shopware orders scanner")

	var cfg mainConfig
	if err := config.Parse("local.env", &cfg); err != nil {
		log.Fatalf("config.Parse: %v", err)
	}

	orderCli, _, err := buildShopwareClients(cfg.Shopware)
	if err != nil {
		log.Fatalf("buildShopwareClients : %v", err)
	}

	midNight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	from := midNight.AddDate(0, 0, -1)
	to := midNight.Add(-time.Nanosecond)

	engine := buildEngine()
	service := orders.NewService(orderCli, engine)
	badOrders, scanned, err := service.ScanOrders(context.Background(), orders.FilterRequest{
		From:                      from,
		To:                        to,
		IncludeCreated:            true,
		IncludeUpdated:            true,
		IncludeDeliveryUpdated:    true,
		IncludeTransactionUpdated: true,
	})
	if err != nil {
		log.Fatalf("failed to scan yesterday orders : %v", err)
	}
	zap.S().Infof("detected %d suspicious orders", len(badOrders))

	renderer := render.NewRenderer("./render/template.twig", cfg.Shopware.BaseURL)
	bytes, err := renderer.RenderHTML(badOrders, scanned)
	if err != nil {
		log.Fatalf("renderer.RenderHTML: %v", err)
	}

	sender := mail.NewSender(cfg.SendGrid)
	err = sender.SendMail(string(bytes))
	if err != nil {
		log.Fatalf("sender.SendMail: %v", err)
	}

	zap.S().Infof("stopping Shopware orders scanner")
}

func buildEngine() rules.Engine {
	rr := map[string]rules.Rule{
		"TRACKING_CODE":       common.ShippedTrackingCode{},
		"PDF_DOCUMENT":        common.ShippedPdfDocument{},
		"RETURN_REFUND_STATE": common.ReturnedRefundedState{},
		"DONE_SHIPPED":        common.DoneDeliveryNotOpen{},
	}

	return rules.NewEngine(rr)
}

func buildShopwareClients(conf config.Shopware) (shopware.OrderService, shopware.ProductService, error) {
	httpCli := resty.New().SetHostURL(conf.BaseURL)

	tokenProvider, err := shopware.NewCredTokenProvider(httpCli, conf.ClientID, conf.ClientSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("shopware.NewCredTokenProvider: %w", err)
	}

	return shopware.NewOrderService(httpCli, tokenProvider),
		shopware.NewProductService(httpCli, tokenProvider), nil
}

func buildLogger() (*zap.Logger, error) {
	ec := zap.NewProductionEncoderConfig()
	ec.LevelKey = "l"
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.RFC3339TimeEncoder
	ec.CallerKey = zapcore.OmitKey
	ec.StacktraceKey = zapcore.OmitKey

	return zap.Config{
		EncoderConfig:    ec,
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}
