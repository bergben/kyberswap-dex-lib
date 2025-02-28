package client

import (
	"context"
	"errors"

	"github.com/KyberNetwork/kutils/klog"
	"github.com/go-resty/resty/v2"

	mxtrading "github.com/KyberNetwork/kyberswap-dex-lib/pkg/liquidity-source/mx-trading"
)

const (
	orderEndpoint = "/order"

	errMsgOrderIsTooSmall = "order is too small"
)

var (
	ErrRFQFailed = errors.New("rfq failed")

	ErrOrderIsTooSmall = errors.New("rfq: order is too small")
)

type HTTPClient struct {
	client *resty.Client
	config *mxtrading.HTTPClientConfig
}

func NewHTTPClient(config *mxtrading.HTTPClientConfig) *HTTPClient {
	client := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout.Duration).
		SetRetryCount(config.RetryCount)

	return &HTTPClient{
		config: config,
		client: client,
	}
}

func (c HTTPClient) Quote(ctx context.Context, params mxtrading.OrderParams) (mxtrading.SignedOrderResult, error) {
	req := c.client.R().SetContext(ctx).SetBody(params)

	var result mxtrading.SignedOrderResult
	var errResult any
	resp, err := req.SetResult(&result).SetError(&errResult).Post(orderEndpoint)
	if err != nil {
		return mxtrading.SignedOrderResult{}, err
	}

	if !resp.IsSuccess() {
		klog.WithFields(ctx, klog.Fields{
			"client":   mxtrading.DexType,
			"response": errResult,
		}).Error("quote failed")
		return mxtrading.SignedOrderResult{}, parseOrderError(errResult)
	}

	return result, nil
}

func parseOrderError(errResult any) error {
	switch errResult {
	case errMsgOrderIsTooSmall:
		return ErrOrderIsTooSmall
	default:
		return ErrRFQFailed
	}
}
