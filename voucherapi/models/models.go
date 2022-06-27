package models

type GetVoucher struct {
	VID         string `json:"VID"`
	UserID      string `json:"UserID"`
	Points      string `json:"Points"`
	Value       string `json:"Value"`
	CreatedDate string `json:"CreatedDate"`
}

type ConsumeVID struct {
	VID          string `json:"VID"`
	UserID       string `json:"UserID"`
	MerchantID   string `json:"MerchantID"`
	Branch       string `json:"Branch"`
	RedeemedDate string `json:"RedeemedDate"`
}

type ReturnMessage struct {
	Ok   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data interface {
	} `json:"data"`
}

type LastBalance struct {
	Balance string `json:"Balance"`
}

type ClaimedFloatFund struct {
	VID       string `json:"VID"`
	ClaimedOn string `json:"ClaimedOn"`
}
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
}

type Admin struct {
	Username string
	Password []byte
}
