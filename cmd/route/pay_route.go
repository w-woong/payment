package route

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/w-woong/common"
	commonport "github.com/w-woong/common/port"
	"github.com/w-woong/payment/delivery"
	"github.com/w-woong/payment/port"
)

func PayRoute(router *mux.Router, conf common.ConfigHttp,
	validator commonport.IDTokenValidators, usc port.PayKcpUsc, userSvc commonport.UserSvc) *delivery.PayKcpHttpHandler {

	handler := delivery.NewPayKcpHttpHandler(time.Duration(conf.Timeout)*time.Second, usc)

	router.HandleFunc("/v1/payment", handler.HandleRegisterCard).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/v1/payment/callback", handler.HandleOrderCardCallback).Methods(http.MethodGet, http.MethodPost)
	// router.HandleFunc("/v1/order/cart", middlewares.AuthIDTokenUserAccountHandler(
	// 	handler.HandleFindByUserID, validator, userSvc,
	// )).Methods(http.MethodGet)
	// router.HandleFunc("/v1/order/cart/_find-or-create", middlewares.AuthIDTokenUserAccountHandler(
	// 	handler.HandleFindOrCreateByUserID, validator, userSvc,
	// )).Methods(http.MethodGet)

	// router.HandleFunc("/test/order/cart", middlewares.AuthIDTokenUserAccountHandler(
	// 	handler.HandleTestRefreshError, validator, userSvc,
	// )).Methods(http.MethodGet)

	// router.HandleFunc("/v1/order/cart/product", middlewares.AuthIDTokenUserAccountHandler(
	// 	handler.HandleAddCartProduct, validator, userSvc,
	// )).Methods(http.MethodPost)

	return handler
}
