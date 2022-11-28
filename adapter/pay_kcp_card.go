package adapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sihttp"
	"github.com/w-woong/payment/dto/kcpdto"
)

type payKcpCard struct {
	client  *sihttp.Client
	baseUrl string
	host    string

	// siteCd           string
	// kcpCertInfo      string
	// allowedPayMethod string
	// returnUrl        string
	// escwUsedYN       string
	// shopName         string
}

// 개발 "https://stg-spl.kcp.co.kr"
// 운영 ?
func NewPayKcpCard(client *http.Client, baseUrl string) *payKcpCard {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Accept-Charset"] = "UTF-8"

	c := sihttp.NewClient(client, sihttp.WithBaseUrl(baseUrl),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
		sihttp.WithDefaultHeaders(headers))

	a := &payKcpCard{
		client:  c,
		baseUrl: baseUrl,
	}
	if u, err := url.Parse(baseUrl); err == nil {
		a.host = u.Host
	}
	return a
}

func (a *payKcpCard) Register(ctx context.Context, req kcpdto.RegistrationRequest) (kcpdto.RegistrationResponse, error) {

	resBody, err := a.client.RequestPostContext(ctx, "/std/tradeReg/register", nil, &req)
	if err != nil {
		return kcpdto.NilRegistrationResponse, err
	}

	res := kcpdto.RegistrationResponse{}
	if err = json.Unmarshal(resBody, &res); err != nil {
		return kcpdto.NilRegistrationResponse, err
	}
	return res, nil
}

func (a *payKcpCard) Approve(ctx context.Context, req kcpdto.ApprovalRequest) (kcpdto.ApprovalResponse, error) {
	// req := kcpdto.ApprovalRequest{
	// 	TranCd:      tranCd,
	// 	SiteCd:      siteCd,
	// 	KcpCertInfo: a.kcpCertInfo,
	// 	EncData:     encData,
	// 	EncInfo:     encInfo,
	// 	OrdrMony:    payAmount,
	// }

	resBody, err := a.client.RequestPostContext(ctx, "/gw/enc/v1/payment", nil, &req)
	if err != nil {
		return kcpdto.NilApprovalResponse, err
	}

	res := kcpdto.ApprovalResponse{}
	if err = json.Unmarshal(resBody, &res); err != nil {
		return kcpdto.NilApprovalResponse, err
	}

	res.SiteCd = req.SiteCd

	return res, nil
}

func (a *payKcpCard) Refund(ctx context.Context, req kcpdto.RefundRequest) (kcpdto.RefundResponse, error) {

	resBody, err := a.client.RequestPostContext(ctx, "/gw/mod/v1/cancel", nil, &req)
	if err != nil {
		return kcpdto.NilRefundResponse, err
	}

	res := kcpdto.RefundResponse{}
	if err = json.Unmarshal(resBody, &res); err != nil {
		return kcpdto.NilRefundResponse, err
	}

	res.SiteCd = req.SiteCd

	return res, nil
}
