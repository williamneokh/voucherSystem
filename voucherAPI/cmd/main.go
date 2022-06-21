package main

import (
	"fmt"
	"github.com/williamneokh/voucherSystem/voucherAPI/config"
	"github.com/williamneokh/voucherSystem/voucherAPI/database"
	"github.com/williamneokh/voucherSystem/voucherAPI/handler"
	"log"
	"net/http"
	"time"
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

	fmt.Println("Listening at port" + portNumber)
	srv := &http.Server{
		Handler: route(),
		Addr:    "localhost" + portNumber,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServeTLS("ssl/localhost.cert.pem", "ssl/localhost.key.pem"))
	//log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(portNumber router))
}
