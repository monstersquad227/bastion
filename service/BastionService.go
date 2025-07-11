package service

import (
	"bastion/model"
	"bastion/repository"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
)

type BastionService struct {
	BastionRepository *repository.BastionRepository
}

func (svc *BastionService) List(userId int) ([]list.Item, error) {
	machines, err := svc.BastionRepository.List(userId)
	if err != nil {
		return nil, err
	}
	items := make([]list.Item, 0, len(machines))
	for _, val := range machines {
		describe := fmt.Sprintf("%s - %s", val.InstanceName, val.PrivateIP)
		items = append(items, model.Item(describe))
	}
	return items, nil
}

func (svc *BastionService) GetPasswordByIp(privateIp string) (string, error) {
	return svc.BastionRepository.GetPassword(privateIp)
}
