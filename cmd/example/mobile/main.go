package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sihttp"
	"github.com/gorilla/mux"
	"github.com/w-woong/payment/dto/kcpdto"
)

var (
	certKey string
	certPem string
)

func init() {
	flag.StringVar(&certKey, "key", "./certs/key.pem", "server key")
	flag.StringVar(&certPem, "pem", "./certs/cert.pem", "server pem")
	flag.Parse()
}
func main() {

	// tmpl := template.Must(template.ParseFiles("resources/html/kcp_mobile.html"))
	tmplTradeRequest := template.Must(template.ParseFiles("resources/html/kcp_mobile_trade_request.html"))
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		req := kcpdto.RegistrationRequest{
			SiteCd:      "T0000",
			KcpCertInfo: "-----BEGIN CERTIFICATE-----MIIDgTCCAmmgAwIBAgIHBy4lYNG7ojANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJLUjEOMAwGA1UECAwFU2VvdWwxEDAOBgNVBAcMB0d1cm8tZ3UxFTATBgNVBAoMDE5ITktDUCBDb3JwLjETMBEGA1UECwwKSVQgQ2VudGVyLjEWMBQGA1UEAwwNc3BsLmtjcC5jby5rcjAeFw0yMTA2MjkwMDM0MzdaFw0yNjA2MjgwMDM0MzdaMHAxCzAJBgNVBAYTAktSMQ4wDAYDVQQIDAVTZW91bDEQMA4GA1UEBwwHR3Vyby1ndTERMA8GA1UECgwITG9jYWxXZWIxETAPBgNVBAsMCERFVlBHV0VCMRkwFwYDVQQDDBAyMDIxMDYyOTEwMDAwMDI0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAppkVQkU4SwNTYbIUaNDVhu2w1uvG4qip0U7h9n90cLfKymIRKDiebLhLIVFctuhTmgY7tkE7yQTNkD+jXHYufQ/qj06ukwf1BtqUVru9mqa7ysU298B6l9v0Fv8h3ztTYvfHEBmpB6AoZDBChMEua7Or/L3C2vYtU/6lWLjBT1xwXVLvNN/7XpQokuWq0rnjSRThcXrDpWMbqYYUt/CL7YHosfBazAXLoN5JvTd1O9C3FPxLxwcIAI9H8SbWIQKhap7JeA/IUP1Vk4K/o3Yiytl6Aqh3U1egHfEdWNqwpaiHPuM/jsDkVzuS9FV4RCdcBEsRPnAWHz10w8CX7e7zdwIDAQABox0wGzAOBgNVHQ8BAf8EBAMCB4AwCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAg9lYy+dM/8Dnz4COc+XIjEwr4FeC9ExnWaaxH6GlWjJbB94O2L26arrjT2hGl9jUzwd+BdvTGdNCpEjOz3KEq8yJhcu5mFxMskLnHNo1lg5qtydIID6eSgew3vm6d7b3O6pYd+NHdHQsuMw5S5z1m+0TbBQkb6A9RKE1md5/Yw+NymDy+c4NaKsbxepw+HtSOnma/R7TErQ/8qVioIthEpwbqyjgIoGzgOdEFsF9mfkt/5k6rR0WX8xzcro5XSB3T+oecMS54j0+nHyoS96/llRLqFDBUfWn5Cay7pJNWXCnw4jIiBsTBa3q95RVRyMEcDgPwugMXPXGBwNoMOOpuQ==-----END CERTIFICATE-----",
			OrdrIdxx:    "TEST1234567890",
			GoodMny:     "1004",
			GoodName:    "TEST상품",
			PayMethod:   "100000000000",
			RetUrl:      "https://192.168.0.92:8099/callback",
			EscwUsed:    "N",
			UserAgent:   "",
		}
		res, err := registerTrade(req)
		if err != nil {
			return
		}
		// data2 := kcpdto.Registration{
		// 	Request:  req,
		// 	Response: res,
		// }
		data := kcpdto.OrderRequest{
			OrdrIdxx: req.OrdrIdxx,
			GoodName: req.GoodName,
			GoodMny:  req.GoodMny,
			BuyrName: "홍길동",
			BuyrTel2: "010-0000-0000",
			BuyrMail: "test@test.co.kr",

			ReqTx:        "pay",
			ShopName:     "TEST SITE",
			SiteCd:       req.SiteCd,
			Currency:     "410",
			EscwUsed:     req.EscwUsed,
			PayMethod:    res.PaymentMethod,
			ActionResult: "",
			VanCode:      "",

			QuotaOpt: "0",

			RetUrl:    req.RetUrl,
			ParamOpt1: "",
			ParamOpt2: "",
			ParamOpt3: "",

			ApprovalKey:   res.ApprovalKey,
			TraceNo:       res.TraceNo,
			PayUrl:        res.PayUrl,
			EncodingTrans: "UTF-8",
		}
		tmplTradeRequest.Execute(w, &data)
	})

	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		_ = dumpRequest(os.Stdout, "ValidateWithJwks", r) // Ignore the error
		if err := r.ParseForm(); err != nil {
			return
		}

		values := r.Form
		res := kcpdto.OrderResponse{
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

		apprRes, err := approve(kcpdto.ApprovalRequest{
			TranCd:      res.TranCd,
			SiteCd:      res.SiteCd,
			KcpCertInfo: "-----BEGIN CERTIFICATE-----MIIDgTCCAmmgAwIBAgIHBy4lYNG7ojANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJLUjEOMAwGA1UECAwFU2VvdWwxEDAOBgNVBAcMB0d1cm8tZ3UxFTATBgNVBAoMDE5ITktDUCBDb3JwLjETMBEGA1UECwwKSVQgQ2VudGVyLjEWMBQGA1UEAwwNc3BsLmtjcC5jby5rcjAeFw0yMTA2MjkwMDM0MzdaFw0yNjA2MjgwMDM0MzdaMHAxCzAJBgNVBAYTAktSMQ4wDAYDVQQIDAVTZW91bDEQMA4GA1UEBwwHR3Vyby1ndTERMA8GA1UECgwITG9jYWxXZWIxETAPBgNVBAsMCERFVlBHV0VCMRkwFwYDVQQDDBAyMDIxMDYyOTEwMDAwMDI0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAppkVQkU4SwNTYbIUaNDVhu2w1uvG4qip0U7h9n90cLfKymIRKDiebLhLIVFctuhTmgY7tkE7yQTNkD+jXHYufQ/qj06ukwf1BtqUVru9mqa7ysU298B6l9v0Fv8h3ztTYvfHEBmpB6AoZDBChMEua7Or/L3C2vYtU/6lWLjBT1xwXVLvNN/7XpQokuWq0rnjSRThcXrDpWMbqYYUt/CL7YHosfBazAXLoN5JvTd1O9C3FPxLxwcIAI9H8SbWIQKhap7JeA/IUP1Vk4K/o3Yiytl6Aqh3U1egHfEdWNqwpaiHPuM/jsDkVzuS9FV4RCdcBEsRPnAWHz10w8CX7e7zdwIDAQABox0wGzAOBgNVHQ8BAf8EBAMCB4AwCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAg9lYy+dM/8Dnz4COc+XIjEwr4FeC9ExnWaaxH6GlWjJbB94O2L26arrjT2hGl9jUzwd+BdvTGdNCpEjOz3KEq8yJhcu5mFxMskLnHNo1lg5qtydIID6eSgew3vm6d7b3O6pYd+NHdHQsuMw5S5z1m+0TbBQkb6A9RKE1md5/Yw+NymDy+c4NaKsbxepw+HtSOnma/R7TErQ/8qVioIthEpwbqyjgIoGzgOdEFsF9mfkt/5k6rR0WX8xzcro5XSB3T+oecMS54j0+nHyoS96/llRLqFDBUfWn5Cay7pJNWXCnw4jIiBsTBa3q95RVRyMEcDgPwugMXPXGBwNoMOOpuQ==-----END CERTIFICATE-----",
			EncData:     res.EncData,
			EncInfo:     res.EncInfo,
			OrdrMony:    res.GoodMny,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(apprRes)

		refund(apprRes)
	})

	// router.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
	// 	_ = dumpRequest(os.Stdout, "ValidateWithJwks", r) // Ignore the error
	// 	if err := r.ParseForm(); err != nil {
	// 		return
	// 	}
	// 	fmt.Println("req_tx", r.Form["req_tx"])
	// 	fmt.Println("req_tx", r.Form["req_tx"])
	// 	fmt.Println("req_tx", r.Form["req_tx"])
	// })

	fileServer := http.FileServer(http.Dir("./resources/"))
	router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", handleFileServe(fileServer)))
	server := http.Server{
		Addr:    ":8099",
		Handler: router,
	}
	// server.ListenAndServe()
	server.ListenAndServeTLS(certPem, certKey)
}

func handleFileServe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

//

func registerTrade(req kcpdto.RegistrationRequest) (kcpdto.RegistrationResponse, error) {
	c := sihttp.DefaultInsecureClient()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Accept-Charset"] = "UTF-8"
	client := sihttp.NewClient(c, sihttp.WithBaseUrl("https://stg-spl.kcp.co.kr"),
		// sihttp.WithRequestOpt(sihttp.WithBearerToken(bearerToken)),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
		sihttp.WithDefaultHeaders(headers),
	)

	// res := TradeRegistrationResponse{}
	// err := client.RequestPostDecodeContext(context.Background(), "/std/tradeReg/register", nil, &req, &res)
	// if err != nil {
	// 	return NilTradeRegistrationResponse, err
	// }
	// fmt.Println(res.String())

	res := kcpdto.RegistrationResponse{}
	resBody, err := client.RequestPostContext(context.Background(), "/std/tradeReg/register", nil, &req)
	if err != nil {
		return kcpdto.NilRegistrationResponse, err
	}
	fmt.Println(string(resBody))
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return kcpdto.NilRegistrationResponse, err
	}
	return res, nil
}

func approve(req kcpdto.ApprovalRequest) (kcpdto.ApprovalResponse, error) {
	c := sihttp.DefaultInsecureClient()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Accept-Charset"] = "UTF-8"
	client := sihttp.NewClient(c, sihttp.WithBaseUrl("https://stg-spl.kcp.co.kr"),
		// sihttp.WithRequestOpt(sihttp.WithBearerToken(bearerToken)),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
		sihttp.WithDefaultHeaders(headers),
	)

	res := kcpdto.ApprovalResponse{}
	resBody, err := client.RequestPostContext(context.Background(), "/gw/enc/v1/payment", nil, &req)
	if err != nil {
		return kcpdto.NilApprovalResponse, err
	}
	fmt.Println(string(resBody))
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return kcpdto.NilApprovalResponse, err
	}
	res.SiteCd = req.SiteCd

	return res, nil
}

func refund(appr kcpdto.ApprovalResponse) (kcpdto.RefundResponse, error) {
	c := sihttp.DefaultInsecureClient()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json; charset=utf-8"
	headers["Accept-Charset"] = "UTF-8"
	client := sihttp.NewClient(c, sihttp.WithBaseUrl("https://stg-spl.kcp.co.kr"),
		// sihttp.WithRequestOpt(sihttp.WithBearerToken(bearerToken)),
		sihttp.WithWriterOpt(sicore.SetJsonEncoder()),
		sihttp.WithReaderOpt(sicore.SetJsonDecoder()),
		sihttp.WithDefaultHeaders(headers),
	)

	pkey, err := loadDecryptedPrivateKey("./certs/dec.pem")
	if err != nil {
		fmt.Println(err)
		return kcpdto.NilRefundResponse, err
	}

	signature, err := sign(pkey, []byte(appr.SiteCd+"^"+appr.Tno+"^"+kcpdto.ModTypeRefundAll)) // 전체취소
	if err != nil {
		fmt.Println(err)
		return kcpdto.NilRefundResponse, err
	}

	req := kcpdto.RefundRequest{
		SiteCd:      appr.SiteCd,
		KcpCertInfo: "-----BEGIN CERTIFICATE-----MIIDgTCCAmmgAwIBAgIHBy4lYNG7ojANBgkqhkiG9w0BAQsFADBzMQswCQYDVQQGEwJLUjEOMAwGA1UECAwFU2VvdWwxEDAOBgNVBAcMB0d1cm8tZ3UxFTATBgNVBAoMDE5ITktDUCBDb3JwLjETMBEGA1UECwwKSVQgQ2VudGVyLjEWMBQGA1UEAwwNc3BsLmtjcC5jby5rcjAeFw0yMTA2MjkwMDM0MzdaFw0yNjA2MjgwMDM0MzdaMHAxCzAJBgNVBAYTAktSMQ4wDAYDVQQIDAVTZW91bDEQMA4GA1UEBwwHR3Vyby1ndTERMA8GA1UECgwITG9jYWxXZWIxETAPBgNVBAsMCERFVlBHV0VCMRkwFwYDVQQDDBAyMDIxMDYyOTEwMDAwMDI0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAppkVQkU4SwNTYbIUaNDVhu2w1uvG4qip0U7h9n90cLfKymIRKDiebLhLIVFctuhTmgY7tkE7yQTNkD+jXHYufQ/qj06ukwf1BtqUVru9mqa7ysU298B6l9v0Fv8h3ztTYvfHEBmpB6AoZDBChMEua7Or/L3C2vYtU/6lWLjBT1xwXVLvNN/7XpQokuWq0rnjSRThcXrDpWMbqYYUt/CL7YHosfBazAXLoN5JvTd1O9C3FPxLxwcIAI9H8SbWIQKhap7JeA/IUP1Vk4K/o3Yiytl6Aqh3U1egHfEdWNqwpaiHPuM/jsDkVzuS9FV4RCdcBEsRPnAWHz10w8CX7e7zdwIDAQABox0wGzAOBgNVHQ8BAf8EBAMCB4AwCQYDVR0TBAIwADANBgkqhkiG9w0BAQsFAAOCAQEAg9lYy+dM/8Dnz4COc+XIjEwr4FeC9ExnWaaxH6GlWjJbB94O2L26arrjT2hGl9jUzwd+BdvTGdNCpEjOz3KEq8yJhcu5mFxMskLnHNo1lg5qtydIID6eSgew3vm6d7b3O6pYd+NHdHQsuMw5S5z1m+0TbBQkb6A9RKE1md5/Yw+NymDy+c4NaKsbxepw+HtSOnma/R7TErQ/8qVioIthEpwbqyjgIoGzgOdEFsF9mfkt/5k6rR0WX8xzcro5XSB3T+oecMS54j0+nHyoS96/llRLqFDBUfWn5Cay7pJNWXCnw4jIiBsTBa3q95RVRyMEcDgPwugMXPXGBwNoMOOpuQ==-----END CERTIFICATE-----",
		KcpSignData: signature,
		ModType:     kcpdto.ModTypeRefundAll,
		Tno:         appr.Tno,
	}

	resBody, err := client.RequestPostContext(context.Background(), "/gw/mod/v1/cancel", nil, &req)
	if err != nil {
		return kcpdto.NilRefundResponse, err
	}
	fmt.Println(string(resBody))

	res := kcpdto.RefundResponse{}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return kcpdto.NilRefundResponse, err
	}

	return res, nil
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

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	writer.Write([]byte("\n"))
	return nil
}
