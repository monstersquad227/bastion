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

type item string

func (i item) FilterValue() string { return "" }

func (svc *BastionService) List() ([]list.Item, error) {
	machines, err := svc.BastionRepository.List(3)
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
