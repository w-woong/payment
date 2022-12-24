package delivery

import (
	"fmt"
	"net/http"
	"time"

	"github.com/w-woong/common/logger"
	"github.com/w-woong/payment/dto/kcpdto"
	"github.com/w-woong/payment/port"
)

type PayKcpHttpHandler struct {
	timeout time.Duration
	usc     port.PayKcpUsc
}

func NewPayKcpHttpHandler(timeout time.Duration, usc port.PayKcpUsc) *PayKcpHttpHandler {
	return &PayKcpHttpHandler{
		timeout: timeout,
		usc:     usc,
	}
}

func (d *PayKcpHttpHandler) HandleRegisterCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderNum := "TEST1234567890"
	var orderAmt float64 = 1004
	productName := "TEST상품"
	buyerName := "홍길동"
	buyerMobile := "010-0000-0000"
	buyerEmail := "test@test.co.kr"
	quota := "0"
	shopName := "TEST SITE"

	orderReq, err := d.usc.Register(ctx, orderNum, orderAmt, productName, buyerName, buyerMobile, buyerEmail, quota, shopName)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if err = d.usc.Order(ctx, w, orderReq); err != nil {
		logger.Error(err.Error())
		return
	}
}

func (d *PayKcpHttpHandler) HandleOrderCardCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	values := r.Form
	orderRes := kcpdto.OrderResponse{
		ResCd:             values.Get("res_cd"),
		ResMsg:            values.Get("res_msg"),
		TranCd:            values.Get("tran_cd"),
		KvpCardCode:       values.Get("kvp_card_code"),
		CardCode:          values.Get("card_code"),
		KcpSelectCardCode: values.Get("kcp_select_card_code"),
		EncData:           values.Get("enc_data"),
		EncInfo:           values.Get("enc_info"),
		CardPayMethod:     values.Get("card_pay_method"),
		CardMaskNo:        values.Get("card_mask_no"),
		CardPointUse:      values.Get("card_point_use"),
		AppPayType:        values.Get("app_pay_type"),
		KcpGroupID:        values.Get("kcp_group_id"),
		EscwUsed:          values.Get("escw_used"),
		ReqTx:             values.Get("req_tx"),
		RtnKeyInfoYn:      values.Get("rtn_key_info_yn"),
		TraceNo:           values.Get("trace_no"),
		BuyrName:          values.Get("buyr_name"),
		BuyrTel1:          values.Get("buyr_tel1"),
		BuyrTel2:          values.Get("buyr_tel2"),
		BuyrMail:          values.Get("buyr_mail"),
		LangFlag:          values.Get("lang_flag"),
		ParamOpt1:         values.Get("param_opt_1"),
		ParamOpt2:         values.Get("param_opt_2"),
		ParamOpt3:         values.Get("param_opt_3"),
		AllyType:          values.Get("ally_type"),
		ConfirmType:       values.Get("confirm_type"),
		GoodName:          values.Get("good_name"),
		DeliTerm:          values.Get("deli_term"),
		RcvrAdd1:          values.Get("rcvr_add1"),
		RcvrAdd2:          values.Get("rcvr_add2"),
		RcvrZipx:          values.Get("rcvr_zipx"),
		RcvrMail:          values.Get("rcvr_mail"),
		BaskCntx:          values.Get("bask_cntx"),
		RetURLMethodType:  values.Get("Ret_URL_Method_Type"),
		PayMethod:         values.Get("pay_method"),
		RcvrTel1:          values.Get("rcvr_tel1"),
		RcvrTel2:          values.Get("rcvr_tel2"),
		SiteCd:            values.Get("site_cd"),
		GoodMny:           values.Get("good_mny"),
		OrdrIdxx:          values.Get("ordr_idxx"),
		EncodingTrans:     values.Get("encoding_trans"),
		RcvrName:          values.Get("rcvr_name"),
		PayMod:            values.Get("pay_mod"),
		GoodInfo:          values.Get("good_info"),
		UsePayMethod:      values.Get("use_pay_method"),
		SmartUseYn:        values.Get("smart_useyn"),
		AllyCode:          values.Get("ally_code"),
		RetURL:            values.Get("Ret_URL"),
	}

	apprRes, err := d.usc.Approve(r.Context(), orderRes)
	if err != nil {
		return
	}

	fmt.Println(apprRes)
}
