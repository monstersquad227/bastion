package service

import (
	"bastion/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func (svc *UserService) GetUserId(account string) (int, error) {
	return svc.UserRepository.GetUserID(account)
}
