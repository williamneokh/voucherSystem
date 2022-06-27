package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/williamneokh/voucherSystem/config"
	"github.com/williamneokh/voucherSystem/database"
	"github.com/williamneokh/voucherSystem/models"
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

type empty struct{}

func Api(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			empty{},
		}
		_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
		return
	}
	_, _ = fmt.Fprintf(w, "You are connected to Voucher System API")
}

//Sponsor allow administrator to add sponsor fund into MasterFund table
func Sponsor(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			empty{},
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
					empty{},
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
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: The characters for sponsor name cannot be more than 36.",
					empty{},
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
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

			//check sponsor amount if it exceeded database allowable characters of 8
			if len(newSponsor.Amount) > 8 {

				w.WriteHeader(http.StatusNotAcceptable)
				unsuccessfulMsg := models.ReturnMessage{
					false,
					"[MS-VOUCHERS]: The characters for sponsor amount cannot be more than 8.",
					empty{},
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
					empty{},
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
				return
			}

			err = newSponsor.DepositMasterFund(params["sponsorid"], newSponsor.SponsorNameOrUserID, newSponsor.Amount)
			if err != nil {
				log.Println(err)
			}
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
			empty{},
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
			empty{},
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
					empty{},
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
					empty{},
				}
				_ = json.NewEncoder(w).Encode(unsuccessfulMsg)

				return
			}

			//Generate Voucher(VID)
			u1 := uuid.NewV4()
			VID := u1.String()

			var isRecordedDatabase = true

			//add 3 var to check the outcome of each go routine
			var isSuccessMasterFund = true
			var isSuccessVoucher = true
			var isSuccessFloatFund = true

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
					isSuccessVoucher = false
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
						empty{},
					}
					_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
					isRecordedDatabase = false
					isSuccessMasterFund = false

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
						empty{},
					}
					_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
					isRecordedDatabase = false
					isSuccessFloatFund = false
				}
			}()

			wg.Wait()
			if isRecordedDatabase == true {
				var getDetails database.DbVoucher
				result := getDetails.GetVoucherDetails(VID)

				w.WriteHeader(http.StatusOK)
				var issueVID models.GetVoucher
				issueVID = models.GetVoucher{VID, gen.UserID, gen.Points, gen.Value, result.CreatedDate}
				successMsg := models.ReturnMessage{
					true,
					"[MS-VOUCHERS]: Generate new voucher, successful",
					issueVID,
				}

				_ = json.NewEncoder(w).Encode(successMsg)
				return
			} else {

				//roll back MasterFund with go routine
				if isSuccessMasterFund == true {
					wg.Add(1)
					go func() {
						err = mFund.RemoveMasterFund(VID, &wg)
						if err != nil {
							log.Println(err)
						}
						fmt.Printf("Successfully row back SponsorIDorVID: %s from MasterFund database\n", VID)
						return
					}()

				}

				//roll back Voucher with go routine
				if isSuccessVoucher == true {
					wg.Add(1)
					go func() {
						err = voucher.RemoveVoucher(VID, &wg)
						if err != nil {
							log.Println(err)
						}
						fmt.Printf("Successfully row back VID: %s from Voucher database\n", VID)
						return
					}()

				}

				//roll back FloatFund with go routine
				if isSuccessFloatFund == true {
					wg.Add(1)
					go func() {
						err = fFund.RemoveFloatFund(VID, &wg)
						if err != nil {
							log.Println(err)
						}
						fmt.Printf("Successfully row back VID: %s from FloatFund database\n", VID)
						return
					}()

				}
				wg.Wait()
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
		if con.VID == "" || con.UserID == "" || con.MerchantID == "" || con.Branch == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				"[MS-VOUCHERS]: Please supply sponsor information in JSON format",
				empty{},
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
				empty{},
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			isRecordDatabase = false
			return

		}
		err = voucher.RedeemVoucher(con.VID, con.MerchantID, con.Branch)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				fmt.Sprintf("[MS-VOUCHERS]: %s", err),
				empty{},
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			isRecordDatabase = false
			return
		}
		if isRecordDatabase == true {
			var getDetails database.DbVoucher
			result := getDetails.GetVoucherDetails(con.VID)
			w.WriteHeader(http.StatusAccepted)

			successConsume := models.ConsumeVID{
				VID:          con.VID,
				UserID:       con.UserID,
				MerchantID:   con.MerchantID,
				Branch:       con.Branch,
				RedeemedDate: result.RedeemedDate,
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

func FundBalance(w http.ResponseWriter, r *http.Request) {
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

	var bal database.DbMasterFund

	latestBalance := bal.FindLatestBalance()
	data := models.LastBalance{Balance: latestBalance}

	successMsg := models.ReturnMessage{
		true,
		"[MS-VOUCHERS]: Pull latest fund balance, successful",
		data,
	}
	_ = json.NewEncoder(w).Encode(successMsg)
	return

}

func MerchantClaims(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		unsuccessfulMsg := models.ReturnMessage{
			false,
			"[MS-VOUCHERS]: Invalid key API_TOKEN",
			empty{},
		}
		_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
		return
	}
	if r.Method == "POST" {

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		var wd database.DbFloatFund
		_ = json.Unmarshal(reqBody, &wd)

		//Validation check if field are empty
		if wd.VID == "" || wd.Branch == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				"[MS-VOUCHERS]: Please supply sponsor information in JSON format",
				empty{},
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			return
		}
		var holdMerchantID = wd.Branch
		err = wd.VendorWithdrawal(wd.VID, wd.Branch)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			unsuccessfulMsg := models.ReturnMessage{
				false,
				fmt.Sprintf("MS-VOUCHERS]: %v", err),
				empty{},
			}
			_ = json.NewEncoder(w).Encode(unsuccessfulMsg)
			return
		}
		var getDetails database.DbFloatFund
		result := getDetails.GetFloatDetails(wd.VID)
		data := models.ClaimedFloatFund{
			VID:       wd.VID,
			ClaimedOn: result.WithdrawalDate,
		}

		successMsg := models.ReturnMessage{
			true,
			fmt.Sprintf("[MS-VOUCHERS]: Successfully claim, fund has been credited to %v's bank", holdMerchantID),
			data,
		}
		_ = json.NewEncoder(w).Encode(successMsg)
		return
	}

}
