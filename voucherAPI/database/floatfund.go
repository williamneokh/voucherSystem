package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type DbFloatFund struct {
	FFund_ID        string `json:"FFund_ID"`
	VID             string `json:"VID"`
	FloatDate       string `json:"FloatDate"`
	FloatValue      string `json:"FloatValue"`
	WithdrawalDate  string `json:"WithdrawalDate"`
	WithdrawalValue string `json:"WithdrawalValue"`
	MerchantID      string `json:"MerchantID"`
	Branch          string `json:"Branch"`
}

//AddFloat transfer from the MasterFund to FloatFund database. FloatFund database is where vendor are allowed to make fund claims.
func (m *DbFloatFund) AddFloat(VID, floatValue string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return err
	}

	defer db.Close()

	query := fmt.Sprintf("INSERT INTO FloatFund (VID, FloatValue, WithdrawalDate, MerchantID, Branch) VALUES('%s','%s','%s','%s','%s')",
		VID, floatValue, "2000-01-01 00:00:00", "OPEN", "OPEN")

	_, err = db.Query(query)
	if err != nil {
		return errors.New("Something went wrong while trying to record VID into FloatFund database")
	}
	return nil
}

//VendorWithdrawal - When vendor submitted their VID claims, VendorWithdrawal match the VID and mark with claimed timestamp and merchantID
//to indicate, fund has been paid to vendor from FloatFund database
func (m *DbFloatFund) VendorWithdrawal(VID, merchantID, branch string) error {

	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return err
	}

	defer db.Close()
	//Validation check if voucher ID exist in database
	results, err := db.Query("SELECT * FROM FloatFund Where VID = ?", VID)
	if err != nil {
		log.Fatal(err)
	}
	var count int
	for results.Next() {
		count++
		err = results.Scan(&m.FFund_ID, &m.VID, &m.FloatDate, &m.FloatValue, &m.WithdrawalDate, &m.MerchantID, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
	}
	if count == 0 {
		return errors.New("voucher ID not found! Please check voucher ID")
	}

	//Validation check if Voucher already been claimed
	if m.MerchantID != "OPEN" {
		return errors.New(fmt.Sprintf("voucher has been claimed before on: %v - by merchant ID: %v", m.WithdrawalDate, m.MerchantID))
	}

	//perform update database with timestamp and merchantID
	now := time.Now()

	query := fmt.Sprintf("UPDATE FloatFund SET WithdrawalDate='%s', MerchantID='%s', Branch='%s' WHERE VID='%s'", now.Format("2006-01-02 15:04:05"), merchantID, branch, VID)

	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil

}
