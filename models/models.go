package models

type GetVoucher struct {
	VID    string `json:"VID"`
	UserID string `json:"UserID"`
	Points string `json:"Points"`
	Value  string `json:"Value"`
}
