package models

type GetVoucher struct {
	VID    string `json:"VID"`
	UserID string `json:"UserID"`
	Points string `json:"Points"`
	Value  string `json:"Value"`
}

type ConsumeVID struct {
	VID        string `json:"VID"`
	UserID     string `json:"UserID"`
	MerchantID string `json:"MerchantID"`
}

type ReturnMessage struct {
	Ok   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data interface {
	} `json:"data"`
}
