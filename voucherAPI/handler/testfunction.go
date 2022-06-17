package handler

import (
	"github.com/williamneokh/voucherSystem/voucherAPI/database"
	"net/http"
)

func TestingGround(w http.ResponseWriter, r *http.Request) {
	if !ValidKey(r) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("401 -Invalid key"))
		return
	}
	// my testing to add voucher into database

	/*var newVoucher database.DbVoucher

	if newVoucher.CheckDuplicatedVID("aaaa") {
		fmt.Fprintf(w, "Duplicated Voucher ID (VID) while trying to insert into Voucher table")
		return
	}

	newVoucher.InsertVoucher("aaaa", "user001", "1000", "5")
	fmt.Fprintf(w, "Add new voucher")

	*/

	//testing withdraw from masterfund

	/* var withdrawMF database.DbMasterFund
	if withdrawMF.CheckSponsorIDorVID("983jdf") {
		fmt.Fprintf(w, "Duplicated Voucher ID(VID) will trying to insert into MasterFund table")
		return
	}
	withdrawMF.WithdrawMasterFund("983jdf", "user001", "5")
	fmt.Fprintf(w, "Test withdrawl from masterfund when inserting new voucher")

	*/

	//test merchant redeem voucher

	/*	var redeem database.DbVoucher
		redeem.RedeemVoucher("aaaa", "mvendor001")

	*/

	var addfloat database.DbFloatFund
	addfloat.AddFloat("039284k", "5")
}
