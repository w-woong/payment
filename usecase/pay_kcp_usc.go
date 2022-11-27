package usecase

import (
	"context"
	"fmt"

	"github.com/w-woong/payment/dto/kcpdto"
	"github.com/w-woong/payment/port"
)

type payKcpUsc struct {
	svc port.PayKcpSvc

	siteCd           string
	kcpCertInfo      string
	allowedPayMethod string
	returnUrl        string
	escwUsedYN       string
}

func NewPayKcpUsc(svc port.PayKcpSvc) *payKcpUsc {
	return &payKcpUsc{
		svc: svc,
	}
}

func (u *payKcpUsc) Register(ctx context.Context,
	orderNum string, orderAmt float64, productName string,
	buyerName, buyerMobile, buyerEmail, quota, shopName string) (kcpdto.OrderRequest, error) {

	goodMny := fmt.Sprintf("%d", int(orderAmt))
	req := kcpdto.RegistrationRequest{
		SiteCd:      u.siteCd,
		KcpCertInfo: u.kcpCertInfo,
		OrdrIdxx:    orderNum,
		GoodMny:     goodMny,
		GoodName:    productName,
		PayMethod:   u.allowedPayMethod,
		RetUrl:      u.returnUrl,
		EscwUsed:    u.escwUsedYN,
		UserAgent:   "",
	}

	res, err := u.svc.Register(ctx, req)
	if err != nil {
		return kcpdto.NilOrderRequest, err
	}

	order := kcpdto.OrderRequest{
		OrdrIdxx: orderNum,
		GoodName: productName,
		GoodMny:  goodMny,
		BuyrName: buyerName,
		BuyrTel2: buyerMobile,
		BuyrMail: buyerEmail,

		ReqTx:        kcpdto.ReqTxPay,
		ShopName:     shopName,
		SiteCd:       u.siteCd,
		Currency:     kcpdto.CurrencyWon,
		EscwUsed:     u.escwUsedYN,
		PayMethod:    kcpdto.PaymentMethodCard,
		ActionResult: kcpdto.ActionResultCard,
		VanCode:      "",

		QuotaOpt: quota,

		RetUrl:    u.returnUrl,
		ParamOpt1: "",
		ParamOpt2: "",
		ParamOpt3: "",

		ApprovalKey:   res.ApprovalKey,
		TraceNo:       res.TraceNo,
		PayUrl:        res.PayUrl,
		EncodingTrans: "UTF-8",
	}

	return order, nil
}

func (u *payKcpUsc) Approve(ctx context.Context, data kcpdto.OrderResponse) (kcpdto.ApprovalResponse, error) {
	req := kcpdto.ApprovalRequest{
		TranCd:      data.TranCd,
		SiteCd:      data.SiteCd,
		KcpCertInfo: u.kcpCertInfo,
		EncData:     data.EncData,
		EncInfo:     data.EncInfo,
		OrdrMony:    data.GoodMny,
	}

	return u.svc.Approve(ctx, req)
}
