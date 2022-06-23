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

type DbVoucher struct {
	Voucher_ID   string `json:"Voucher_ID"`
	VID          string `json:"VID"`
	UserID       string `json:"UserID"`
	UserPoints   string `json:"UserPoints"`
	CreatedDate  string `json:"CreatedDate"`
	VoucherValue string `json:"VoucherValue"`
	RedeemedDate string `json:"RedeemedDate"`
	MerchantID   string `json:"MerchantID"`
	Branch       string `json:"Branch"`
}

//RedeemVoucher is when merchant consumed user's voucher.
func (m *DbVoucher) RedeemVoucher(VID, merchant, branch string) error {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	now := time.Now()

	query := fmt.Sprintf("UPDATE Voucher SET RedeemedDate='%s', MerchantID='%s', Branch='%s' WHERE VID='%s'", now.Format("2006-01-02 15:04:05"), merchant, branch, VID)

	_, err = db.Query(query)
	if err != nil {
		return err

	}
	return nil
}

//InsertVoucher is when user making a trade of their points for voucher
func (m *DbVoucher) InsertVoucher(VID, userID, userPoints, voucherValue string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return err
	}

	defer db.Close()

	query := fmt.Sprintf("INSERT INTO Voucher (VID, UserID, UserPoints, VoucherValue, RedeemedDate, MerchantID, Branch) VALUES('%s','%s','%s','%s','%s','%s','%s')",
		VID, userID, userPoints, voucherValue, "2000-01-01 00:00:00", "OPEN", "OPEN")

	_, err = db.Query(query)
	if err != nil {
		return errors.New("Something went wrong while trying to record VID into Voucher database")
	}
	return nil

}

//RemoveVoucher to remove voucher record

func (m DbVoucher) RemoveVoucher(VID string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return err
	}

	defer db.Close()

	query := fmt.Sprintf("DELETE FROM Voucher where VID ='%s'", VID)
	_, err = db.Query(query)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to remove VID: %s from database", VID))
	}
	return nil
}

//CheckDuplicatedVID check if the VID is duplicated in Voucher database
func (m *DbVoucher) CheckDuplicatedVID(vid string) bool {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	results, err := db.Query("SELECT * FROM Voucher Where VID = ?", vid)
	if err != nil {
		log.Fatal(err)
	}
	for results.Next() {
		err = results.Scan(&m.Voucher_ID, &m.VID, &m.UserID, &m.UserPoints, &m.CreatedDate, &m.VoucherValue, &m.RedeemedDate, &m.MerchantID, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
		if m.VID == "" {
			return false
		} else {
			return true
		}
	}
	return false
}

//ValidateVoucher to check the authenticity of a voucher
func (m *DbVoucher) ValidateVoucher(vid, userId string) error {

	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM Voucher Where VID = ?", vid)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(vid, userId)
	i := 0
	for results.Next() {
		i++
		err = results.Scan(&m.Voucher_ID, &m.VID, &m.UserID, &m.UserPoints, &m.CreatedDate, &m.VoucherValue, &m.RedeemedDate, &m.MerchantID, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}

	}

	//check user's voucher is available in database
	if i == 0 {
		return errors.New("Voucher cannot be found in database")
	} else {
		//check user is it the rightful owner of the voucher
		if m.UserID != userId {
			return errors.New("UserID doesnt match the correct user in the database")
		} else {
			//check if the voucher has be use before
			if m.MerchantID == "OPEN" {
				//fmt.Println("Valid unused voucher")
			} else {
				return errors.New("voucher has been used")
			}
		}

	}
	return nil
}

func (m *DbVoucher) TotalVoucherIssued() (int, []DbVoucher) {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	//Find out the latest balance from database
	results, err := db.Query("SELECT * FROM Voucher ORDER BY Voucher_ID DESC LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	var totalValue int
	var created []DbVoucher
	for results.Next() {
		err = results.Scan(&m.Voucher_ID, &m.VID, &m.UserID, &m.UserPoints, &m.CreatedDate, &m.VoucherValue, &m.RedeemedDate, &m.MerchantID, &m.Branch)
		if err != nil {
			log.Fatal(err)
		}
		value, _ := strconv.Atoi(m.VoucherValue)
		totalValue += value
		data := DbVoucher{
			Voucher_ID:   m.Voucher_ID,
			VID:          m.VID,
			UserID:       m.UserID,
			UserPoints:   m.UserPoints,
			CreatedDate:  m.CreatedDate,
			VoucherValue: m.VoucherValue,
			RedeemedDate: m.RedeemedDate,
			MerchantID:   m.MerchantID,
			Branch:       m.Branch,
		}
		created = append(created, data)
	}
	return totalValue, created
}
