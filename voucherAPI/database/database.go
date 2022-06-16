package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/williamneokh/voucherSystem/voucherAPI/config"
	"log"
	"strconv"
)

type DbMasterFund struct {
	Mfund_ID            int    `json:"Mfund_ID"`
	TransactionType     string `json:"TransactionType"`
	SponsorIDOrVID      string `json:"SponsorIDOrVID"`
	SponsorNameOrUserID string `json:"SponsorNameOrUserID"`
	TransactionDate     string `json:"TransactionDate"`
	Amount              string `json:"Amount"`
	BalancedFund        string `json:"BalancedFund"`
}

var vip *config.Config

func NewDataRepo(a *config.Config) {
	vip = a
}

//InsertFund function create new incoming sponsorship fund into the mySql table MasterFund
func (m *DbMasterFund) InsertFund(transactionType, sponsorIDorVID, sponsorNameOrUserID, amount string) {

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
		//fmt.Println(m.BalancedFund)
	}
	//
	//Convert m.BalancedFund and deposit into INT
	lastBalance, err := strconv.Atoi(m.BalancedFund)
	newDeposit, err := strconv.Atoi(amount)
	//
	//Return sum from adding new deposit with latest balance from database
	var sum = lastBalance + newDeposit

	//convert back to string to be use for
	newBalanced := strconv.Itoa(sum)

	fmt.Println(newBalanced)
	//
	query := fmt.Sprintf("INSERT INTO MasterFund (TransactionType, SponsorIDOrVID, "+
		"SponsorNameOrUserID, Amount, BalancedFund) VALUES('%s','%s','%s','%s','%s')",
		transactionType, sponsorIDorVID, sponsorNameOrUserID, amount, newBalanced)

	_, err = db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
}

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
