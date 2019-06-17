package models

type RegisterRsp struct {
	Uuid     string `json:"uuid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	PassWord string `json:"password"`
}
