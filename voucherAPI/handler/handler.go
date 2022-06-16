package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/williamneokh/voucherSystem/voucherAPI/database"
	"io/ioutil"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are connected to Voucher System API")
}

func Sponsor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if r.Header.Get("Content-type") == "application/json" {

		if r.Method == "POST" {

			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			var newSponsor database.DbMasterFund

			json.Unmarshal(reqBody, &newSponsor)

			//check for any empty field
			if newSponsor.TransactionType == "" || params["sponsorid"] == "" || newSponsor.SponsorNameOrUserID == "" || newSponsor.Amount == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, _ = w.Write([]byte("422 - Please supply sponsor information in JSON format"))
				return
			}

			//check for any wrong type of data
			if newSponsor.TransactionType != "Deposit" {
				w.WriteHeader(http.StatusNotAcceptable)
				_, _ = w.Write([]byte("406 - Transaction type for sponsor funding must be name exactly \"Deposit\""))
				return
			}
			//check if the characters use has been exceeded the database size
			if len(params["sponsorid"]) > 8 {
				w.WriteHeader(http.StatusNotAcceptable)
				_, _ = w.Write([]byte("406 - The characters cannot be more than 8."))
				return
			}
			//check if the sponsorID has been use.
			if newSponsor.CheckSponsorIDorVID(params["sponsorid"]) {

				w.WriteHeader(http.StatusNotAcceptable)
				_, _ = w.Write([]byte("406 - The Sponsor ID or VID has been used, please create new sponser ID or Voucher."))
				return
			}

			newSponsor.InsertFund(newSponsor.TransactionType, params["sponsorid"], newSponsor.SponsorNameOrUserID, newSponsor.Amount)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("201 - Funds added: " + params["sponsorid"]))
		}

	}
}
