package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/williamneokh/voucherSystem/voucherAPI/config"
	"github.com/williamneokh/voucherSystem/voucherAPI/database"
	"github.com/williamneokh/voucherSystem/voucherAPI/handler"
	"log"
	"net/http"
)

const portNumber = ":3000"

func main() {
	vip, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	database.ViperDatabase(&vip)
	handler.ViperHandler(&vip)

	//sponsor = make(map[string]handler.Sponsorship)

	router := mux.NewRouter()
	router.HandleFunc("/api", handler.Home)
	router.HandleFunc("/api/sponsor/{sponsorid}", handler.Sponsor).Methods("POST")
	router.HandleFunc("/api/masterfund", handler.AllMasterFundRecords).Methods("POST")
	router.HandleFunc("/api/getvoucher", handler.TestingGround).Methods("POST")

	fmt.Println("Listening at port" + portNumber)

	log.Fatal(http.ListenAndServe(portNumber, router))
}
