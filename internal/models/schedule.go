package models

type Schedule struct {
	Minute []int `json:"minute"`
	Hour   []int `json:"hour"`
	Mday   []int `json:"mday"`
	Wday   []int `json:"wday"`
	Month  []int `json:"month"`
}
