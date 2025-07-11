package utils

import (
	"bastion/model"
	"bastion/service"
)

func ValidateUser(user, password string, svc service.BastionServiceInterface) bool {
	items, err := svc.List()
	if err != nil {
		return false
	}
	model.DefaultVMList = items
	ok := ComparePassword(user, password)
	if !ok {
		return false
	}
	return true
}
