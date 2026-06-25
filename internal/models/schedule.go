package models

type Schedule struct {
	Minute []int `json:"minute"`
	Hour   []int `json:"hour"`
	Mday   []int `json:"mday"`
	Month  []int `json:"month"`
	Wday   []int `json:"wday"`
}
