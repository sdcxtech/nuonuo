package nuonuo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type TokenController interface {
	GetToken(ctx context.Context) (string, error)
}

type permanentToken struct {
	token string
}

func NewPermanentToken(token string) TokenController {
	return &permanentToken{
		token: token,
	}
}

func (pt *permanentToken) GetToken(ctx context.Context) (string, error) {
	return pt.token, nil
}

type oauthToken struct {
	appKey    string
	appSecret string

	restyClient *resty.Client

	token       string
	isPermanent bool
	expiresTime time.Time
	mu          sync.Mutex
}

func NewOAuthToken(appKey, appSecret string) TokenController {
	return &oauthToken{
		appKey:      appKey,
		appSecret:   appSecret,
		restyClient: resty.New(),
	}
}

func (ot *oauthToken) GetToken(ctx context.Context) (string, error) {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if !ot.isValid() {
		if err := ot.refreshToken(ctx); err != nil {
			return "", err
		}
	}

	return ot.token, nil
}

func (ot *oauthToken) isValid() bool {
	return ot.token != "" && (ot.isPermanent || time.Now().Before(ot.expiresTime))
}

func (ot *oauthToken) refreshToken(ctx context.Context) error {
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	resp, err := ot.restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8").
		SetFormData(map[string]string{
			"client_id":     ot.appKey,
			"client_secret": ot.appSecret,
			"grant_type":    "client_credentials",
		}).
		ForceContentType("application/json").
		SetResult(&result).
		Post("https://open.nuonuo.com/accessToken")
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("http status: %s, body: %s", resp.Status(), resp.Body())
	}

	ot.token = result.AccessToken

	if result.ExpiresIn < 0 {
		ot.isPermanent = true
	} else {
		ot.expiresTime = time.Now().Add(time.Second * time.Duration(result.ExpiresIn-5))
	}

	return nil
}
