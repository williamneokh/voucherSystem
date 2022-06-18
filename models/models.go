package models

type GetVoucher struct {
	VID    string `json:"VID"`
	UserID string `json:"UserID"`
	Points string `json:"Points"`
	Value  string `json:"Value"`
}

type ConsumeVID struct {
	Status     string `json:"Status"`
	Message    string `json:"Message"`
	VID        string `json:"VID"`
	UserID     string `json:"UserID"`
	MerchantID string `json:"MerchantID"`
}
