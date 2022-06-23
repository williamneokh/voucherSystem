package handler

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/williamneokh/voucherSystem/database"
	"github.com/williamneokh/voucherSystem/models"
	"github.com/williamneokh/voucherSystem/render"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
	var vValue = make(map[string]int)
	totalVoucherValue, created := v.TotalVoucherIssued()

	record["voucher"] = created
	vValue["total"] = totalVoucherValue

	var bal database.DbMasterFund
	latestBalance := bal.FindLatestBalance()

	fund := make(map[string]string)
	fund["balance"] = latestBalance
	render.RenderTemplate(w, "dashboard.page.tmpl", &models.TemplateData{
		StringMap: fund,
		IntMap:    vValue,
		Data:      record,
	})
}
