package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
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

//AddFloat transfer from the MasterFund to FloatFund database. FloatFund database is where vendor are allowed to make fund claims.
func (m *DbFloatFund) AddFloat(VID, floatValue string, group *sync.WaitGroup) error {
	defer group.Done()
	db, err := sql.Open(vip.DBDriver, vip.DBSource)

	if err != nil {
		return err
	}

	defer db.Close()

	query := fmt.Sprintf("INSERT INTO FloatFund (VID, FloatValue, WithdrawalDate, MerchantID) VALUES('%s','%s','%s','%s')",
		VID, floatValue, "2000-01-01 00:00:00", "OPEN")

	_, err = db.Query(query)
	if err != nil {
		return errors.New("Something went wrong while trying to record VID into FloatFund database")
	}
	return nil
}

//VendorWithdrawal - When vendor submitted their VID claims, VendorWithdrawal match the VID and mark with claimed timestamp and merchantID
//to indicate, fund has been paid to vendor from FloatFund database
func (m *DbFloatFund) VendorWithdrawal(VID, merchantID string) {

	//Match VID against database

	//if no match found, return error, VID no found match found

	//else perform update database with timestamp and merchantID

}
