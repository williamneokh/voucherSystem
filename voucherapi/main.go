package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/williamneokh/voucherSystem/config"
	"github.com/williamneokh/voucherSystem/database"
	"github.com/williamneokh/voucherSystem/handler"
	"github.com/williamneokh/voucherSystem/render"
	"log"
	"net/http"
)

var app config.AppConfig

func route() http.Handler {
	router := mux.NewRouter()

	//Admin Dashboard frontend
	router.HandleFunc("/", handler.Home)
	router.HandleFunc("/login", handler.Login)
	router.HandleFunc("/logout", handler.Logout)
	router.HandleFunc("/dashboard", handler.Dashboard)
	router.HandleFunc("/addfund", handler.AddFund)
	router.HandleFunc("/deposit", handler.DepositMasterFund)

	//API Backend connection
	router.HandleFunc("/api", handler.Api)
	//router.HandleFunc("/api/sponsor/{sponsorid}", handler.Sponsor).Methods("POST")
	router.HandleFunc("/api/masterfund", handler.AllMasterFundRecords).Methods("POST")
	router.HandleFunc("/api/fundbalance", handler.FundBalance).Methods("GET")
	router.HandleFunc("/api/getvoucher", handler.GetVoucher).Methods("POST")
	router.HandleFunc("/api/consumevid", handler.ConsumeVID).Methods("POST")
	router.HandleFunc("/api/merchantclaims", handler.MerchantClaims).Methods("POST")

	//Serve file
	fileServer := http.FileServer(http.Dir("./images/"))
	//router.Handle("/images/*", http.StripPrefix("/images", fileServer))

	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fileServer))

	return router

}

const portNumber = ":3000"

func main() {
	handler.Initial()
	app.InProduction = false
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UserCache = false

	render.NewTemplates(&app)
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
	}
	log.Fatal(srv.ListenAndServeTLS("ssl/localhost.cert.pem", "ssl/localhost.key.pem"))
	//log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(portNumber router))
}
