package token

import "time"

type Token struct {
	AccessToken string
	ExpiresIn   time.Duration
}

type Refresher interface {
	RefreshToken(appKey, appSecret string) (*Token, error)
}
