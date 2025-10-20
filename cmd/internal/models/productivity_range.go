package models

type ProductivityRange struct {
	Ok       int `json:"ok"`
	Alert    int `json:"alert"`
	Critical int `json:"critical"`
}
