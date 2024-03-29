package nuonuo

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_QueryInvoiceRedConfirm(t *testing.T) {
	c := newClient()

	resp, err := c.QueryInvoiceRedConfirm(
		context.Background(),
		&QueryInvoiceRedConfirmRequest{
			Identity: "0",
			BillID:   "non-exist",
		},
	)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Total)
}

func newClient() *Client {
	return New(
		"https://sdk.nuonuo.com/open/v1/services",
		os.Getenv("NUONUO_APP_KEY"),
		os.Getenv("NUONUO_APP_SECRET"),
		os.Getenv("NUONUO_TAX_ID"),
		NewPermanentToken(os.Getenv("NUONUO_TOKEN")),
	)
}
