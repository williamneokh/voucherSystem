package database

import (
	"database/sql"
	"fmt"
	"log"
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
}

func (m *DbVoucher) RedeemVoucher(VID, merchant string) {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	now := time.Now()

	query := fmt.Sprintf("UPDATE Voucher SET RedeemedDate='%s', MerchantID='%s' WHERE VID='%s'", now.Format("2006-01-02 15:04:05"), merchant, VID)

	_, err = db.Query(query)
	if err != nil {
		log.Fatal(err)

	}
}

func (m *DbVoucher) InsertVoucher(VID, userID, userPoints, voucherValue string) {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	query := fmt.Sprintf("INSERT INTO Voucher (VID, UserID, UserPoints, VoucherValue,"+
		"RedeemedDate, MerchantID) VALUES('%s','%s','%s','%s','%s','%s')",
		VID, userID, userPoints, voucherValue, "2000-01-01 00:00:00", "OPEN")

	_, err = db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
}

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
		err = results.Scan(&m.Voucher_ID, &m.VID, &m.UserID, &m.UserPoints, &m.CreatedDate, &m.VoucherValue, &m.RedeemedDate, &m.MerchantID)
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
