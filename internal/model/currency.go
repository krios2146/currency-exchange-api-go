package model

type Currency struct {
	Id       int64  `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"name"`
	Sign     string `json:"sign"`
}
