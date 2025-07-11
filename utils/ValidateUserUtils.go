package utils

import (
	"bastion/model"
	"bastion/service"
	"fmt"
)

func ValidateUser(user, password string, bastionSvc service.BastionServiceInterface, userSvc service.UserServiceInterface) bool {
	// 验证用户名密码是否正确
	if err := service.LdapClient.Bind(
		fmt.Sprintf("cn=%s,ou=%s,dc=%s,dc=%s", user, "technology", "chengdd", "com"), password); err != nil {
		return false
	}
	
	// 获取UserID
	userId, err := userSvc.GetUserId(user)
	if err != nil {
		return false
	}

	// 获取MachineList
	items, err := bastionSvc.List(userId)
	if err != nil {
		return false
	}

	model.DefaultVMList = items

	return true
}
