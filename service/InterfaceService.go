package service

import "github.com/charmbracelet/bubbles/list"

type BastionServiceInterface interface {
	List(userId int) ([]list.Item, error)
}

type UserServiceInterface interface {
	GetUserId(account string) (int, error)
}
