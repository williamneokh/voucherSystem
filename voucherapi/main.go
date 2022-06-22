package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/williamneokh/voucherSystem/config"
	"github.com/williamneokh/voucherSystem/database"
	"github.com/williamneokh/voucherSystem/handler"
	"log"
	"net/http"
	"time"
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

const portNumber = ":3000"

func main() {
	vip, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	database.ViperDatabase(&vip)
	handler.ViperHandler(&vip)

	//sponsor = make(map[string]handler.Sponsorship)

	fmt.Println("Listening at port" + portNumber)
	srv := &http.Server{
		Handler: route(),
		Addr:    portNumber,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServeTLS("ssl/localhost.cert.pem", "ssl/localhost.key.pem"))
	//log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(portNumber router))
}
