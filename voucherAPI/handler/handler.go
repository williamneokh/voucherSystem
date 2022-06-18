package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/williamneokh/voucherSystem/models"
	"github.com/williamneokh/voucherSystem/voucherAPI/config"
	"github.com/williamneokh/voucherSystem/voucherAPI/database"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var vip *config.Config

func ViperHandler(a *config.Config) {
	vip = a
}

//ValidKey check if client key matches api key
func ValidKey(r *http.Request) bool {
	v := r.Header.Get("Key")

	if v == vip.ApiToken {
		return true
	} else {
		return false
	}

}

//Home just a return a message if you are connected to the api
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are connected to Voucher System API")
}

//Sponsor allow administrator to add sponsor fund into MasterFund table
func Sponsor(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("401 -Invalid key"))
		return
	}

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
			if params["sponsorid"] == "" || newSponsor.SponsorNameOrUserID == "" || newSponsor.Amount == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, _ = w.Write([]byte("422 - Please supply sponsor information in JSON format"))
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
				_, _ = w.Write([]byte("406 - The Sponsor ID or VID has been used, please create new sponsor ID or Voucher."))
				return
			}

			newSponsor.DepositMasterFund(params["sponsorid"], newSponsor.SponsorNameOrUserID, newSponsor.Amount)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("201 - Funds added: " + params["sponsorid"]))
		}

	}
}

//AllMasterFundRecords allow administrator to list all the date in Masterfund table
func AllMasterFundRecords(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("401 -Invalid key"))
		return
	}

	var MasterFund database.DbMasterFund

	MasterFund.ListMasterFundRecords(w)

	// returns all the courses in JSON

}

func GetVoucher(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("401 -Invalid key"))
		return
	}
	if r.Header.Get("Content-type") == "application/json" {
		var mFund database.DbMasterFund
		if r.Method == "POST" {

			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			var gen models.GetVoucher

			json.Unmarshal(reqBody, &gen)

			if gen.UserID == "" || gen.Points == "" || gen.Value == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, _ = w.Write([]byte("422 - Please supply sponsor information in JSON format"))
				return
			}
			//check Masterfund if there is enough fund to generate voucher
			if !mFund.CheckMasterFund(gen.Value) {
				w.WriteHeader(http.StatusPaymentRequired)
				_, _ = w.Write([]byte("402 - insufficient balance in MasterFund"))
				return
			}

			//Generate Voucher(VID)
			u1 := uuid.NewV4()
			VID := u1.String()

			var wg sync.WaitGroup

			var voucher database.DbVoucher
			wg.Add(3)
			//Spawn go routine for recording into Voucher database
			go func() {
				err = voucher.InsertVoucher(VID, gen.UserID, gen.Points, gen.Value, &wg)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprintf("500 - InternalServerError: %v", err)))
				}
			}()

			//Spawn go routine for recording into MasterFund database and update the fund balance
			go func() {
				err = mFund.WithdrawMasterFund(VID, gen.UserID, gen.Value, &wg)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprintf("500 - InternalServerError: %v", err)))
				}
			}()

			var fFund database.DbFloatFund
			//Spawn go routine for recording into FloatFund database
			go func() {
				err = fFund.AddFloat(VID, gen.Value, &wg)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(fmt.Sprintf("500 - InternalServerError: %v", err)))
				}
			}()

			wg.Wait()
			w.WriteHeader(http.StatusOK)
			var success models.GetVoucher
			success = models.GetVoucher{VID, gen.UserID, gen.Points, gen.Value}

			_ = json.NewEncoder(w).Encode(success)
		}
	}
}

func ConsumeVID(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("401 -Invalid key"))
		return
	}
	if r.Method == "POST" {

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var con models.ConsumeVID

		json.Unmarshal(reqBody, &con)
		if con.VID == "" || con.UserID == "" || con.MerchantID == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("422 - Please supply sponsor information in JSON format"))
			return
		}

		var voucher database.DbVoucher
		err = voucher.ValidateVoucher(con.VID, con.UserID)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(fmt.Sprintf("%s", err)))

		}
		err = voucher.RedeemVoucher(con.VID, con.MerchantID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("%s", err)))
		}
		w.WriteHeader(http.StatusAccepted)
		//_, _ = w.Write([]byte(fmt.Sprintf("202 - Successfuly consumed\nVID: %s\nFor user: %s\nTo Merchant: %s\n", con.VID, con.UserID, con.MerchantID)))
		successMsg := models.ConsumeVID{
			Status:     "202",
			Message:    "Successfully Consumed",
			VID:        con.VID,
			UserID:     con.UserID,
			MerchantID: con.MerchantID,
		}
		_ = json.NewEncoder(w).Encode(successMsg)
	}
}
