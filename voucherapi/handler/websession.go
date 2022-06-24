package handler

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/williamneokh/voucherSystem/database"
	"github.com/williamneokh/voucherSystem/models"
	"github.com/williamneokh/voucherSystem/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

var mapAdmin = map[string]models.Admin{}
var mapSessions = map[string]string{}

func Initial() {
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost) //should not appear after go action 1
	mapAdmin["admin"] = models.Admin{"admin", bPassword}
}

func AlreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := mapSessions[myCookie.Value]
	_, ok := mapAdmin[username]
	return ok
}

//Home just a return a message if you are connected to the api
func Home(w http.ResponseWriter, r *http.Request) {
	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})

}

func Login(w http.ResponseWriter, r *http.Request) {
	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	// process form submission
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username, password)
		// check if user exist with username
		myAdmin, ok := mapAdmin[username]
		if !ok {
			http.Error(w, "Username do not match", http.StatusUnauthorized)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myAdmin.Password, []byte(password))
		if err != nil {
			http.Error(w, "password do not match", http.StatusUnauthorized)
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(w, myCookie)
		mapSessions[myCookie.Value] = username
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := r.Cookie("myCookie")
	// delete the session
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, myCookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	var v database.DbVoucher
	var record = make(map[string]interface{})

	_, created := v.TotalVoucherIssued()
	_, used := v.TotalVoucherUsed()

	record["created"] = created
	record["used"] = used

	var bal database.DbMasterFund
	latestBalance := bal.FindLatestBalance()

	fund := make(map[string]string)
	fund["balance"] = latestBalance
	render.RenderTemplate(w, "dashboard.page.tmpl", &models.TemplateData{
		StringMap: fund,
		Data:      record,
	})
}

//AddFund allow administrator to add sponsor fund into MasterFund table
func AddFund(w http.ResponseWriter, r *http.Request) {

	if !AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	//var f database.DbMasterFund
	//
	//f.DepositMasterFund()
	var bal database.DbMasterFund
	latestBalance := bal.FindLatestBalance()

	fund := make(map[string]string)
	fund["balance"] = latestBalance

	render.RenderTemplate(w, "addfund.page.tmpl", &models.TemplateData{
		StringMap: fund,
	})
}

func DepositMasterFund(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if r.Method == http.MethodPost {
		SponsorIDOrVID := r.FormValue("SponsorIDOrVID")
		SponsorNameOrUserID := r.FormValue("SponsorNameOrUserID")
		Amount := r.FormValue("Amount")

		//validation check for form input is empty
		if SponsorNameOrUserID == "" || SponsorIDOrVID == "" || Amount == "" {
			http.Error(w, "form input not found", http.StatusNotAcceptable)

			return
		}
		//validation check sponsor id if it exceeded database allowable characters of 8
		if len(SponsorIDOrVID) > 8 {
			http.Error(w, "Sponsor ID cannot be more than 8 varchar", http.StatusNotAcceptable)
			return
		}
		//validation check sponsor name if it exceeded database allowable characters of 36
		if len(SponsorNameOrUserID) > 36 {
			http.Error(w, "Sponsor name cannot be more than 36 varchar", http.StatusNotAcceptable)
			return
		}

		//validation check amount input is integer value
		_, err := strconv.Atoi(Amount)
		if err != nil {

			http.Error(w, fmt.Sprintf("expecting integer value but got : %v", Amount), http.StatusNotAcceptable)
			return
		}
		//validation check sponsor amount if it exceeded database allowable characters of 8
		if len(Amount) > 8 {
			http.Error(w, "the characters use for amount cannot be more than 8 varchar", http.StatusNotAcceptable)
			return
		}

		//validation check if the sponsorID has been used before
		var ns database.DbMasterFund

		if ns.CheckSponsorIDorVID(SponsorIDOrVID) {
			http.Error(w, "the Sponsor ID has been used before, please give another unique Sponsor ID and resubmit again", http.StatusNotAcceptable)
			return
		}
		ns.DepositMasterFund(SponsorIDOrVID, SponsorNameOrUserID, Amount)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

}
