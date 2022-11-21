package refresher

import (
	"time"

	"github.com/sdcxtech/nuonuo/token"
)

type permanentToken struct {
	token     string
	expiresIn time.Duration
}

func NewPermanentToken(token string, expiresIn time.Duration) token.Refresher {
	return &permanentToken{
		token:     token,
		expiresIn: expiresIn,
	}
}

func (pt *permanentToken) RefreshToken(appKey, appSecret string) (*token.Token, error) {
	return &token.Token{
		AccessToken: pt.token,
		ExpiresIn:   pt.expiresIn,
	}, nil
}
