package models

type LoginRsp struct {
	Token    string `json:"token"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"`
}
