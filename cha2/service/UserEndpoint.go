package service

type UserRequest struct {
	UserId int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"data"`
}
