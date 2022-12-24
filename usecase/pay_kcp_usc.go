package usecase

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"text/template"

	"github.com/w-woong/payment/dto/kcpdto"
	"github.com/w-woong/payment/port"
)

type payKcpUsc struct {
	svc port.PayKcpSvc

	siteCd               string
	kcpCertInfo          string
	allowedPayMethod     string
	returnUrl            string
	escwUsedYN           string
	privateKeyPath       string
	privateKey           *rsa.PrivateKey
	tradeRequestHtmlFile string
	tradeRequestTemplate *template.Template
}

func NewPayKcpUsc(svc port.PayKcpSvc, siteCd, kcpCertInfo, allowedPayMethod, returnUrl string,
	privateKeyPath string, tradeRequestHtmlFile string) (*payKcpUsc, error) {

	pkey, err := loadDecryptedPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	tmpl := template.Must(template.ParseFiles(tradeRequestHtmlFile))

	return &payKcpUsc{
		svc:                  svc,
		siteCd:               siteCd,
		kcpCertInfo:          kcpCertInfo,
		allowedPayMethod:     allowedPayMethod,
		returnUrl:            returnUrl,
		escwUsedYN:           "N",
		privateKeyPath:       privateKeyPath,
		privateKey:           pkey,
		tradeRequestHtmlFile: tradeRequestHtmlFile,
		tradeRequestTemplate: tmpl,
	}, nil
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

func (u *payKcpUsc) Order(ctx context.Context, w io.Writer, order kcpdto.OrderRequest) error {
	return u.tradeRequestTemplate.Execute(w, &order)
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

func (u *payKcpUsc) RefundAll(ctx context.Context, siteCd, tno string) (kcpdto.RefundResponse, error) {
	modType := kcpdto.ModTypeRefundAll

	signature, err := sign(u.privateKey, []byte(siteCd+"^"+tno+"^"+modType))
	if err != nil {
		return kcpdto.NilRefundResponse, err
	}
	req := kcpdto.RefundRequest{
		SiteCd:      siteCd,
		KcpCertInfo: u.kcpCertInfo,
		KcpSignData: signature,
		ModType:     modType,
		Tno:         tno,
	}

	return u.svc.Refund(ctx, req)
}

func loadDecryptedPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	privatePem, _ := pem.Decode(bytes)
	fmt.Println(privatePem)

	key, err := x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
func sign(priavateKey *rsa.PrivateKey, data []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priavateKey, crypto.SHA256, h.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}
