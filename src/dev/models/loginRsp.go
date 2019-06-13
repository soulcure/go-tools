package models

type LoginRsp struct {
	Id       int64  `json:"id"`
	Token    string `json:"token"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"`
}
