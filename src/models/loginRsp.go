package models

type LoginRsp struct {
	Token string `json:"token"`
	Uuid  string `json:"uuid"`
}
