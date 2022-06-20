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
	"strconv"
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
	_, _ = fmt.Fprintf(w, "You are connected to Voucher System API")
}

//Sponsor allow administrator to add sponsor fund into MasterFund table
func Sponsor(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			"",
		}
		_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
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

			_ = json.Unmarshal(reqBody, &newSponsor)

			//check for any empty field
			if newSponsor.SponsorNameOrUserID == "" || newSponsor.Amount == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: Please supply sponsor information in JSON format",
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
				return
			}

			//check sponsor id if it exceeded database allowable characters of 8
			if len(params["sponsorid"]) > 8 {
				w.WriteHeader(http.StatusNotAcceptable)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: The characters for sponsor id cannot be more than 8.",
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
				return
			}
			//check sponsor name if it exceeded database allowable characters of 36
			if len(newSponsor.SponsorNameOrUserID) > 36 {
				w.WriteHeader(http.StatusNotAcceptable)
				_, _ = w.Write([]byte("406 - The characters for sponsor name cannot be more than 36."))
				return
			}
			//Make sure sponsor amount is a valid integer
			_, err = strconv.Atoi(newSponsor.Amount)
			if err != nil {
				w.WriteHeader(http.StatusNotAcceptable)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					fmt.Sprintf("[MS-VOUCHERS]: Expecting integer only but got: '%v'", newSponsor.Amount),
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)

				return
			}

			//check sponsor amount if it exceeded database allowable characters of 7
			if len(newSponsor.Amount) > 8 {

				w.WriteHeader(http.StatusNotAcceptable)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: The characters for sponsor amount cannot be more than 8.",
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
				return
			}

			//check if the sponsorID has been use.
			if newSponsor.CheckSponsorIDorVID(params["sponsorid"]) {

				w.WriteHeader(http.StatusNotAcceptable)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: The Sponsor ID or VID has been used, please create new sponsor ID or Voucher",
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
				return
			}

			newSponsor.DepositMasterFund(params["sponsorid"], newSponsor.SponsorNameOrUserID, newSponsor.Amount)
			w.WriteHeader(http.StatusCreated)
			successMsg := models.ReturnMessage{
				true,
				"[MS-VOUCHERS]: Sponsor fund deposit, successful",
				fmt.Sprintf("Code: %v", params["sponsorid"]),
			}

			_ = json.NewEncoder(w).Encode(successMsg)
		}

	}
}

//AllMasterFundRecords allow administrator to list all the date in Masterfund table
func AllMasterFundRecords(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			"",
		}
		_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
		return
	}

	var MasterFund database.DbMasterFund

	MasterFund.ListMasterFundRecords(w)

	// returns all the courses in JSON

}

func GetVoucher(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			"",
		}
		_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
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

			_ = json.Unmarshal(reqBody, &gen)

			//Validation check if field are empty
			if gen.UserID == "" || gen.Points == "" || gen.Value == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: Please supply sponsor information in JSON format",
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)

				return
			}
			//Validation check if there are enough fund to generate new voucher
			if !mFund.CheckMasterFund(gen.Value) {
				w.WriteHeader(http.StatusPaymentRequired)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: Insufficient balance in MasterFund",
					"",
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)

				return
			}

			//Generate Voucher(VID)
			u1 := uuid.NewV4()
			VID := u1.String()

			var isRecordedDatabase = true

			var wg sync.WaitGroup

			var voucher database.DbVoucher
			wg.Add(3)
			//Spawn go routine for recording into Voucher database
			go func() {
				err = voucher.InsertVoucher(VID, gen.UserID, gen.Points, gen.Value, &wg)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					unsuccessfulMsg := models.ReturnMessage{
						false,
						fmt.Sprintf("[MS-VOUCHERS]: InternalServerError: %v", err),
						"",
					}
					_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
					isRecordedDatabase = false
				}
			}()

			//Spawn go routine for recording into MasterFund database and update the fund balance
			go func() {
				err = mFund.WithdrawMasterFund(VID, gen.UserID, gen.Value, &wg)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					unsuccessfulMsg := models.ReturnMessage{
						false,
						fmt.Sprintf("[MS-VOUCHERS]: InternalServerError: %v", err),
						"",
					}
					_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
					isRecordedDatabase = false
				}
			}()

			var fFund database.DbFloatFund
			//Spawn go routine for recording into FloatFund database
			go func() {
				err = fFund.AddFloat(VID, gen.Value, &wg)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					unsuccessfulMsg := models.ReturnMessage{
						false,
						fmt.Sprintf("[MS-VOUCHERS]: InternalServerError: %v", err),
						"",
					}
					_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
					isRecordedDatabase = false
				}
			}()

			wg.Wait()
			if isRecordedDatabase == true {
				w.WriteHeader(http.StatusOK)
				var issueVID models.GetVoucher
				issueVID = models.GetVoucher{VID, gen.UserID, gen.Points, gen.Value}
				successMsg := models.ReturnMessage{
					true,
					"[MS-VOUCHERS]: Generate new voucher, successful",
					issueVID,
				}

				_ = json.NewEncoder(w).Encode(successMsg)
				return
			}
		}
	}
}

func ConsumeVID(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			"",
		}
		_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
		return
	}
	if r.Method == "POST" {

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var con models.ConsumeVID

		_ = json.Unmarshal(reqBody, &con)
		if con.VID == "" || con.UserID == "" || con.MerchantID == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				"[MS-VOUCHERS]: Please supply sponsor information in JSON format",
				"",
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			return
		}
		var isRecordDatabase = true
		var voucher database.DbVoucher
		err = voucher.ValidateVoucher(con.VID, con.UserID)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				fmt.Sprintf("[MS-VOUCHERS]: %s", err),
				"",
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			isRecordDatabase = false
			return

		}
		err = voucher.RedeemVoucher(con.VID, con.MerchantID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				fmt.Sprintf("[MS-VOUCHERS]: %s", err),
				"",
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			isRecordDatabase = false
			return
		}
		if isRecordDatabase == true {
			w.WriteHeader(http.StatusAccepted)

			successConsume := models.ConsumeVID{
				VID:        con.VID,
				UserID:     con.UserID,
				MerchantID: con.MerchantID,
			}
			successMsg := models.ReturnMessage{
				true,
				"[MS-VOUCHERS]: Consume voucher, successful",
				successConsume,
			}
			_ = json.NewEncoder(w).Encode(successMsg)
			return
		}
	}
}
