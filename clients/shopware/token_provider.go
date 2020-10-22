package shopware

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"sync"
	"time"
)

type TokenProvider interface {
	GetToken() string
}

func NewCredTokenProvider(client *resty.Client, clientID, clientSecret string) (TokenProvider, error) {
	c := &credTokenProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		client:       client,
	}

	if _, err := c.updateToken(context.Background()); err != nil {
		return nil, fmt.Errorf("updateToken failed: %w", err)
	}

	go c.scheduleUpdateToken()

	return c, nil
}

type credTokenProvider struct {
	clientID     string
	clientSecret string
	client       *resty.Client

	activeAccessToken string //guarded by mutex
	mu                sync.RWMutex
}

func (c *credTokenProvider) GetToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.activeAccessToken
}

func (c *credTokenProvider) scheduleUpdateToken() {
	for {
		tokenResp, err := c.updateToken(context.Background())
		if err != nil {
			zap.S().Errorf("updateToken : %v", err)
			time.Sleep(time.Second)
			continue
		}
		// 0.9 factor should refrain from requests being failed due to expired token
		time.Sleep(time.Duration(tokenResp.ExpiresIn*9/10) * time.Second)
	}
}

func (c *credTokenProvider) updateToken(ctx context.Context) (tokenResponse, error) {
	request := tokenRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		GrantType:    "client_credentials",
	}
	path := "/api/oauth/token"

	var result tokenResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post(path)

	if err := checkHttpResp(resp, err); err != nil {
		return tokenResponse{}, err
	}

	c.mu.Lock()
	c.activeAccessToken = result.AccessToken
	c.mu.Unlock()

	return result, nil
}
