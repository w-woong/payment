package port

import (
	"context"

	"github.com/w-woong/payment/dto/kcpdto"
)

type PayKcpSvc interface {
	Register(ctx context.Context, req kcpdto.RegistrationRequest) (kcpdto.RegistrationResponse, error)
	Approve(ctx context.Context, req kcpdto.ApprovalRequest) (kcpdto.ApprovalResponse, error)
}

type PayKcpUsc interface {
	Register(ctx context.Context,
		orderNum string, orderAmt float64, productName string,
		buyerName, buyerMobile, buyerEmail, quota, shopName string) (kcpdto.OrderRequest, error)

	Approve(ctx context.Context, data kcpdto.OrderResponse) (kcpdto.ApprovalResponse, error)
}
