package service

import "errors"

type IUserService interface {
	GetName(userId int) string
	DelUser(userId int) error
}

type UserService struct {
}

func (UserService) GetName(userId int) string {
	if userId == 101 {
		return "shanghai"
	}

	return "wangsi"
}

func (UserService) DelUser(userId int) error {
	if userId == 101 {
		return errors.New("denied access")
	}
	return nil
}
