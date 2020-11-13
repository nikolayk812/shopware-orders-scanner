package main

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/nikolayk812/shopware-orders-scanner/checks"
	"github.com/nikolayk812/shopware-orders-scanner/checks/common"
	"github.com/nikolayk812/shopware-orders-scanner/clients/shopware"
	"github.com/nikolayk812/shopware-orders-scanner/config"
	"github.com/nikolayk812/shopware-orders-scanner/consumers/html"
	"github.com/nikolayk812/shopware-orders-scanner/consumers/mail"
	"github.com/nikolayk812/shopware-orders-scanner/orders"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
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
	defer zap.S().Infof("stopping Shopware orders scanner")

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

	if cfg.SendGrid.Enabled {
		sender := mail.NewSender(cfg.Shopware, cfg.SendGrid)
		_, err = sender.Consume(badOrders, scanned)
		if err != nil {
			log.Fatalf("sender.Consume : %v", err)
		}
		return
	}

	htmlRenderer := html.NewRenderer("./consumers/html/template.twig", cfg.Shopware.BaseURL)
	document, err := htmlRenderer.Consume(badOrders, scanned)
	if err != nil {
		log.Fatalf("htmlRenderer.Consume : %v", err)
	}
	fileName := "./reports/" + time.Now().Format("report_01-02-2006_15:04") + ".html"
	if err := ioutil.WriteFile(fileName, document.Bytes, 0644); err != nil {
		log.Fatalf("WriteFile : %v", err)
	}
}

func buildEngine() checks.Engine {
	rr := map[string]checks.Check{
		"TRACKING_CODE":       common.ShippedTrackingCode{},
		"PDF_DOCUMENT":        common.ShippedPdfDocument{},
		"RETURN_REFUND_STATE": common.ReturnedRefundedState{},
		"DONE_SHIPPED":        common.DoneDeliveryNotOpen{},
	}

	return checks.NewEngine(rr)
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
