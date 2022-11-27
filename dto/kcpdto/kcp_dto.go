package kcpdto

import "encoding/json"

var (
	NilRegistrationResponse = RegistrationResponse{}
	NilOrderRequest         = OrderRequest{}
)

const (
	CurrencyWon = "410"
	CurrencyUSD = "840"
)

const (
	ReqTxPay = "pay"
)
const (
	PaymentMethodCard = "CARD"
)
const (
	ActionResultCard = "card"
)

type Registration struct {
	Request  RegistrationRequest
	Response RegistrationResponse
}
type RegistrationRequest struct {
	SiteCd      string `json:"site_cd"`
	KcpCertInfo string `json:"kcp_cert_info"`
	OrdrIdxx    string `json:"ordr_idxx"`
	GoodMny     string `json:"good_mny"`
	GoodName    string `json:"good_name"`

	// 결제 수단 정보 설정
	// 결제에 필요한 결제 수단 정보를 설정합니다.
	//
	// 신용카드 : 100000000000, 계좌이체 : 010000000000, 가상계좌 : 001000000000
	// 포인트   : 000100000000, 휴대폰   : 000010000000, 상품권   : 000000001000
	//
	// 위와 같이 설정한 경우 표준웹에서 설정한 결제수단이 표시됩니다.
	// 표준웹에서 여러 결제수단을 표시하고 싶으신 경우 설정하시려는 결제
	// 수단에 해당하는 위치에 해당하는 값을 1로 변경하여 주십시오.
	//
	// 예) 신용카드, 계좌이체, 가상계좌를 동시에 표시하고자 하는 경우
	// pay_method = "111000000000"
	// 신용카드(100000000000), 계좌이체(010000000000), 가상계좌(001000000000)에
	// 해당하는 값을 모두 더해주면 됩니다.
	//
	// ※ 필수
	//  KCP에 신청된 결제수단으로만 결제가 가능합니다.
	PayMethod string `json:"pay_method"`

	RetUrl    string `json:"Ret_URL"`
	EscwUsed  string `json:"escw_used"`
	UserAgent string `json:"user_agent"`
}

type RegistrationResponse struct {
	Code        string `json:"Code"`
	Message     string `json:"Message"`
	ApprovalKey string `json:"approvalKey"`
	TraceNo     string `json:"traceNo"`
	PayUrl      string `json:"PayUrl"`

	PaymentMethod string `json:"paymentMethod"`
	RequestURI    string `json:"request_URI"`
}

func (o *RegistrationResponse) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}

// 주문요청
type OrderRequest struct {
	OrdrIdxx string `json:"ordr_idxx"`
	GoodName string `json:"good_name"`
	GoodMny  string `json:"good_mny"`
	BuyrName string `json:"buyr_name"`
	BuyrTel2 string `json:"buyr_tel2"` // mobile
	BuyrMail string `json:"buyr_mail"`

	ReqTx        string `json:"req_tx"` // fixed, "pay"
	ShopName     string `json:"shop_name"`
	SiteCd       string `json:"site_cd"`
	Currency     string `json:"currency"`  // 원화: 410, USD: 840
	EscwUsed     string `json:"escw_used"` // fixed, "N"
	PayMethod    string `json:"pay_method"`
	ActionResult string `json:"ActionResult"`
	VanCode      string `json:"van_code"`

	// options
	QuotaOpt string `json:"quotaopt"` // 할부 : 0 ~ 12

	RetUrl    string `json:"Ret_URL"`
	ParamOpt1 string `json:"param_opt_1"`
	ParamOpt2 string `json:"param_opt_2"`
	ParamOpt3 string `json:"param_opt_3"`

	ApprovalKey string `json:"approvalKey"`
	TraceNo     string `json:"traceNo"`
	PayUrl      string `json:"PayUrl"`

	EncodingTrans string `json:"encoding_trans"` // fixed, "UTF-8"
}

func (o *OrderRequest) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}

// 인증완료후 돌아오는 데이터, 이 데이터를 이용해 승인 요청
type OrderResponse struct {
	ResCd             string `json:"res_cd"`
	ResMsg            string `json:"res_msg"`
	TranCd            string `json:"tran_cd"`
	KvpCardCode       string `json:"kvp_card_code"`
	CardCode          string `json:"card_code"`
	KcpSelectCardCode string `json:"kcp_select_card_code"`
	EncData           string `json:"enc_data"`
	EncInfo           string `json:"enc_info"`
	CardPayMethod     string `json:"card_pay_method"`
	CardMaskNo        string `json:"card_mask_no"`
	CardPointUse      string `json:"card_point_use"`
	AppPayType        string `json:"app_pay_type"`
	KcpGroupID        string `json:"kcp_group_id"`
	EscwUsed          string `json:"escw_used"`
	ReqTx             string `json:"req_tx"`
	RtnKeyInfoYn      string `json:"rtn_key_info_yn"`
	TraceNo           string `json:"trace_no"`
	BuyrName          string `json:"buyr_name"`
	BuyrTel1          string `json:"buyr_tel1"`
	BuyrTel2          string `json:"buyr_tel2"`
	BuyrMail          string `json:"buyr_mail"`
	LangFlag          string `json:"lang_flag"`
	ParamOpt1         string `json:"param_opt_1"`
	ParamOpt2         string `json:"param_opt_2"`
	ParamOpt3         string `json:"param_opt_3"`
	AllyType          string `json:"ally_type"`
	ConfirmType       string `json:"confirm_type"`
	GoodName          string `json:"good_name"`
	DeliTerm          string `json:"deli_term"`
	RcvrAdd1          string `json:"rcvr_add1"`
	RcvrAdd2          string `json:"rcvr_add2"`
	RcvrZipx          string `json:"rcvr_zipx"`
	RcvrMail          string `json:"rcvr_mail"`
	BaskCntx          string `json:"bask_cntx"`
	RetURLMethodType  string `json:"Ret_URL_Method_Type"`
	PayMethod         string `json:"pay_method"`
	RcvrTel1          string `json:"rcvr_tel1"`
	RcvrTel2          string `json:"rcvr_tel2"`
	SiteCd            string `json:"site_cd"`
	GoodMny           string `json:"good_mny"`
	OrdrIdxx          string `json:"ordr_idxx"`
	EncodingTrans     string `json:"encoding_trans"`
	RcvrName          string `json:"rcvr_name"`
	PayMod            string `json:"pay_mod"`
	GoodInfo          string `json:"good_info"`
	UsePayMethod      string `json:"use_pay_method"`
	SmartUseYn        string `json:"smart_useyn"`
	AllyCode          string `json:"ally_code"`
	RetURL            string `json:"Ret_URL"`
}

func (o *OrderResponse) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}

var (
	NilApprovalResponse = ApprovalResponse{}
)

type ApprovalRequest struct {
	TranCd      string `json:"tran_cd"`
	SiteCd      string `json:"site_cd"`
	KcpCertInfo string `json:"kcp_cert_info"`
	EncData     string `json:"enc_data"`
	EncInfo     string `json:"enc_info"`
	OrdrMony    string `json:"ordr_mony"`
}
type ApprovalResponse struct {
	ResCd    string `json:"res_cd"`
	ResMsg   string `json:"res_msg"`
	ResEnMsg string `json:"res_en_msg"`
	Tno      string `json:"tno"`
	Amount   string `json:"amount"`

	// card
	OrderNo            string `json:"order_no"`
	MallTaxNo          string `json:"mall_taxno"`
	PartcancYN         string `json:"partcanc_yn"`
	Noinf              string `json:"noinf"`
	CouponMny          string `json:"coupon_mny"`
	IspIssuerCd        string `json:"isp_issuer_cd"`
	PgTxid             string `json:"pg_txid"`
	CardBinType01      string `json:"card_bin_type_01"`
	TraceNo            string `json:"trace_no"`
	CardMny            string `json:"card_mny"`
	ResVatMny          string `json:"res_vat_mny"`
	CaOrderID          string `json:"ca_order_id"`
	ResTaxFlag         string `json:"res_tax_flag"`
	AcquName           string `json:"acqu_name"`
	CardNo             string `json:"card_no"`
	Quota              string `json:"quota"`
	VanCd              string `json:"van_cd"`
	IspPartnerCd       string `json:"isp_partner_cd"`
	AcquCd             string `json:"acqu_cd"`
	CertNo             string `json:"cert_no"`
	VanApptime         string `json:"van_apptime"`
	IspIssuerNm        string `json:"isp_issuer_nm"`
	ResFreeMny         string `json:"res_free_mny"`
	PayMethod          string `json:"pay_method"`
	CardBinBankCd      string `json:"card_bin_bank_cd"`
	BizxNumb           string `json:"bizx_numb"`
	EscwYN             string `json:"escw_yn"`
	JoinCd             string `json:"join_cd"`
	AppTime            string `json:"app_time"`
	CardBinType02      string `json:"card_bin_type_02"`
	CardCd             string `json:"card_cd"`
	CardName           string `json:"card_name"`
	MchtTaxno          string `json:"mcht_taxno"`
	ResGreenDepositMny string `json:"res_green_deposit_mny"`
	ResTaxMny          string `json:"res_tax_mny"`
	AppNo              string `json:"app_no"`
}
