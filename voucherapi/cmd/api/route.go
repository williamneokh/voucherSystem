package main

import (
	"github.com/gorilla/mux"
	"github.com/williamneokh/voucherSystem/handler"
	"net/http"
)

func route() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api", handler.Home)
	router.HandleFunc("/api/sponsor/{sponsorid}", handler.Sponsor).Methods("POST")
	router.HandleFunc("/api/masterfund", handler.AllMasterFundRecords).Methods("POST")
	router.HandleFunc("/api/fundbalance", handler.FundBalance).Methods("GET")
	router.HandleFunc("/api/getvoucher", handler.GetVoucher).Methods("POST")
	router.HandleFunc("/api/consumevid", handler.ConsumeVID).Methods("POST")
	router.HandleFunc("/api/merchantclaims", handler.MerchantClaims).Methods("POST")

	return router

}
