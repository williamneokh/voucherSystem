package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type DbFloatFund struct {
	FFund_ID       string `json:"FFund_ID"`
	VID            string `json:"VID"`
	FloatDate      string `json:"FloatDate"`
	FloatValue     string `json:"FloatValue"`
	WithdrawalDate string `json:"WithdrawalDate"`
	Branch         string `json:"Branch"`
}

//AddFloat transfer from the MasterFund to FloatFund database. FloatFund database is where vendor are allowed to make fund claims.
func (m *DbFloatFund) AddFloat(VID, floatValue string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return errors.New("failed to connect to FloatFund database")
	}

	defer db.Close()

	query := fmt.Sprintf("INSERT INTO FloatFund (VID, FloatValue, WithdrawalDate, Branch) VALUES('%s','%s','%s','%s')",
		VID, floatValue, "2000-01-01 00:00:00", "OPEN")

	_, err = db.Query(query)
	if err != nil {
		return errors.New("Something went wrong while trying to record VID into FloatFund database")
	}
	return nil
}

//RemoveFloatFund delete FloatFund record
func (m *DbFloatFund) RemoveFloatFund(VID string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM FloatFund WHERE VID='%s'", VID)

	_, err = db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to delete VID: %s from floatfund database", VID))
	}
	return nil
}

//VendorWithdrawal - When vendor submitted their VID claims, VendorWithdrawal match the VID and mark with claimed timestamp and merchantID
//to indicate, fund has been paid to vendor from FloatFund database
func (m *DbFloatFund) VendorWithdrawal(VID, branch string) error {

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
		err = results.Scan(&m.FFund_ID, &m.VID, &m.FloatDate, &m.FloatValue, &m.WithdrawalDate, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
	}
	if count == 0 {
		return errors.New("voucher ID not found! Please check voucher ID")
	}

	//Validation check if Voucher already been claimed
	if m.Branch != "OPEN" {
		return errors.New(fmt.Sprintf("voucher has been claimed before on: %v - by merchant ID: %v", m.WithdrawalDate, m.Branch))
	}

	//Validation if the branch consume is equal to the same branch that claim the voucher
	var details DbVoucher
	result := details.GetVoucherDetails(VID)
	if branch != result.Branch {
		return errors.New(fmt.Sprintf("voucher was used on: %v at location: %v, a different branch: %v trying to claim the voucher.", result.RedeemedDate, result.Branch, branch))
	}

	//perform update database with timestamp and merchantID
	now := time.Now()

	query := fmt.Sprintf("UPDATE FloatFund SET WithdrawalDate='%s', Branch='%s' WHERE VID='%s'", now.Format("2006-01-02 15:04:05"), branch, VID)

	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil

}

//GetFloatDetails search vid and retrieve information of that data
func (m *DbFloatFund) GetFloatDetails(VID string) DbFloatFund {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM FloatFund Where VID = ?", VID)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		err = results.Scan(&m.FFund_ID, &m.VID, &m.FloatDate, &m.FloatValue, &m.WithdrawalDate, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
		data := DbFloatFund{
			FFund_ID:       m.FFund_ID,
			VID:            m.VID,
			FloatDate:      m.FloatDate,
			FloatValue:     m.FloatValue,
			WithdrawalDate: m.WithdrawalDate,
			Branch:         m.Branch,
		}
		return data

	}
	return DbFloatFund{}
}

//FloatFundBalance calculate the available floating fund balance
func (m *DbFloatFund) FloatFundBalance() int {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//Find out the latest balance from database
	results, err := db.Query("SELECT * FROM FloatFund WHERE Branch = 'OPEN'")
	if err != nil {
		log.Fatal(err)
	}
	var totalValue int

	for results.Next() {
		err = results.Scan(&m.FFund_ID, &m.VID, &m.FloatDate, &m.FloatValue, &m.WithdrawalDate, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
		value, _ := strconv.Atoi(m.FloatValue)
		totalValue += value

	}
	return totalValue
}

//MerchantClaimed calculate the total amount issued to merchant
func (m *DbFloatFund) MerchantClaimed() int {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//Find out the latest balance from database
	results, err := db.Query("SELECT * FROM FloatFund WHERE Branch != 'OPEN'")
	if err != nil {
		log.Fatal(err)
	}
	var totalValue int

	for results.Next() {
		err = results.Scan(&m.FFund_ID, &m.VID, &m.FloatDate, &m.FloatValue, &m.WithdrawalDate, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
		value, _ := strconv.Atoi(m.FloatValue)
		totalValue += value

	}
	return totalValue
}
