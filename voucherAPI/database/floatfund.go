package database

import (
	"database/sql"
	"fmt"
	"log"
)

type DbFloatFund struct {
	FFund_ID        string `json:"FFund_ID"`
	VID             string `json:"VID"`
	FloatDate       string `json:"FloatDate"`
	FloatValue      string `json:"FloatValue"`
	WithdrawalDate  string `json:"WithdrawalDate"`
	WithdrawalValue string `json:"WithdrawalValue"`
	Merchant_ID     string `json:"Merchant_ID"`
}

func (m *DbFloatFund) AddFloat(VID, floatValue string) {
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	query := fmt.Sprintf("INSERT INTO FloatFund (VID, FloatValue, WithdrawalDate, MerchantID) VALUES('%s','%s','%s','%s')",
		VID, floatValue, "2000-01-01 00:00:00", "OPEN")

	_, err = db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *DbFloatFund) VendorWithdrawal() {

}
