package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/williamneokh/voucherSystem/config"
	"github.com/williamneokh/voucherSystem/models"
	"log"
	"net/http"
	"strconv"
	"sync"
)

//DbMasterFund struct is associated with MasterFund table
type DbMasterFund struct {
	Mfund_ID            int    `json:"Mfund_ID"`
	TransactionType     string `json:"TransactionType"`
	SponsorIDOrVID      string `json:"SponsorIDOrVID"`
	SponsorNameOrUserID string `json:"SponsorNameOrUserID"`
	TransactionDate     string `json:"TransactionDate"`
	Amount              string `json:"Amount"`
	BalancedFund        string `json:"BalancedFund"`
}

//MasterFund is a map that stored MasterFund database into local memory.
var MasterFund map[string]DbMasterFund

var vip *config.Config

//ViperDatabase load the viper configuration so that the shared env variable can be use in database package
func ViperDatabase(a *config.Config) {
	vip = a
}

//DepositMasterFund function create new incoming sponsorship fund into the mySql table MasterFund and update new balance
func (m *DbMasterFund) DepositMasterFund(sponsorID, sponsorName, amount string) {

	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//Find out the latest balance from database
	results, err := db.Query("SELECT * FROM MasterFund ORDER BY MFund_ID DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		err = results.Scan(&m.Mfund_ID, &m.TransactionType, &m.SponsorIDOrVID, &m.SponsorNameOrUserID, &m.TransactionDate, &m.Amount, &m.BalancedFund)
		if err != nil {
			log.Fatal(err)
		}

	}

	//Convert m.BalancedFund and deposit into INT
	lastBalance, err := strconv.Atoi(m.BalancedFund)
	newDeposit, err := strconv.Atoi(amount)

	//Return sum from adding new deposit with latest balance from database
	var sum = lastBalance + newDeposit

	//convert back to string to be use for
	newBalanced := strconv.Itoa(sum)

	fmt.Println(newBalanced)

	query := fmt.Sprintf("INSERT INTO MasterFund (TransactionType, SponsorIDOrVID, "+
		"SponsorNameOrUserID, Amount, BalancedFund) VALUES('%s','%s','%s','%s','%s')",
		"Deposit", sponsorID, sponsorName, amount, newBalanced)

	_, err = db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
}

//CheckSponsorIDorVID this function is part of validation check duplicated SponsorID or voucher ID(VID)
func (m *DbMasterFund) CheckSponsorIDorVID(id string) bool {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM MasterFund Where SponsorIDOrVID = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		err = results.Scan(&m.Mfund_ID, &m.TransactionType, &m.SponsorIDOrVID, &m.SponsorNameOrUserID, &m.TransactionDate, &m.Amount, &m.BalancedFund)
		if err != nil {
			log.Fatal(err)
		}
		if m.SponsorIDOrVID == "" {
			return false
		} else {
			return true
		}
	}
	return false
}

//ListMasterFundRecords function returned all records of MasterFund from database into JSON Format
func (m *DbMasterFund) ListMasterFundRecords(w http.ResponseWriter) {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM MasterFund")
	if err != nil {
		log.Fatal(err)
	}
	var newRecord DbMasterFund
	var count int
	for results.Next() {
		count++
		err = results.Scan(&m.Mfund_ID, &m.TransactionType, &m.SponsorIDOrVID, &m.SponsorNameOrUserID, &m.TransactionDate, &m.Amount, &m.BalancedFund)
		if err != nil {
			log.Fatal(err)
		}
		newRecord = DbMasterFund{m.Mfund_ID, m.TransactionType, m.SponsorIDOrVID, m.SponsorNameOrUserID, m.TransactionDate, m.Amount, m.BalancedFund}
		successMsg := models.ReturnMessage{
			true,
			fmt.Sprintf("[MS-VOUCHERS]: Result: %d", count),
			newRecord,
		}

		_ = json.NewEncoder(w).Encode(successMsg)

	}
}

//WithdrawMasterFund insert voucher into Masterfund database, it also update the masterfund balance
func (m *DbMasterFund) WithdrawMasterFund(VID, userID, amount string, group *sync.WaitGroup) error {
	defer group.Done()

	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return errors.New("failed to connect to MaterFund database")
	}

	defer db.Close()

	//Find out the latest balance from database
	results, err := db.Query("SELECT * FROM MasterFund ORDER BY MFund_ID DESC LIMIT 1")
	if err != nil {
		return err
	}
	for results.Next() {
		err = results.Scan(&m.Mfund_ID, &m.TransactionType, &m.SponsorIDOrVID, &m.SponsorNameOrUserID, &m.TransactionDate, &m.Amount, &m.BalancedFund)
		if err != nil {
			return err
		}

	}

	//Convert m.BalancedFund and deposit into INT
	lastBalance, err := strconv.Atoi(m.BalancedFund)
	newVoucherValue, err := strconv.Atoi(amount)

	//Return sum from adding new deposit with latest balance from database
	var sum = lastBalance - newVoucherValue

	//convert back to string to be use for
	newBalanced := strconv.Itoa(sum)

	query := fmt.Sprintf("INSERT INTO MasterFund (TransactionType, SponsorIDOrVID, "+
		"SponsorNameOrUserID, Amount, BalancedFund) VALUES('%s','%s','%s','%s','%s')",
		"Withdrawal", VID, userID, amount, newBalanced)

	_, err = db.Query(query)
	if err != nil {
		return errors.New("Something went wrong while trying to record VID into MasterFund database")
	}
	return nil
}

//RemoveMasterFund - delete record from database
func (m *DbMasterFund) RemoveMasterFund(VID string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM MasterFund where SponsorIDOrVID='%s'", VID)

	_, err = db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to remove voucher ID: %s from MasterFund database", VID))
	}
	return nil
}

func (m *DbMasterFund) CheckMasterFund(amount string) bool {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//Find out the latest balance from database
	results, err := db.Query("SELECT * FROM MasterFund ORDER BY MFund_ID DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		err = results.Scan(&m.Mfund_ID, &m.TransactionType, &m.SponsorIDOrVID, &m.SponsorNameOrUserID, &m.TransactionDate, &m.Amount, &m.BalancedFund)
		if err != nil {
			log.Fatal(err)
		}

	}

	//Convert m.BalancedFund and deposit into INT
	lastBalance, err := strconv.Atoi(m.BalancedFund)
	amountInt, err := strconv.Atoi(amount)

	if lastBalance > amountInt {
		return true
	} else {
		return false
	}
}

func (m *DbMasterFund) FindLatestBalance() string {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM MasterFund ORDER BY MFund_ID DESC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		err = results.Scan(&m.Mfund_ID, &m.TransactionType, &m.SponsorIDOrVID, &m.SponsorNameOrUserID, &m.TransactionDate, &m.Amount, &m.BalancedFund)
		if err != nil {
			log.Fatal(err)
		}

	}
	return m.BalancedFund

}
