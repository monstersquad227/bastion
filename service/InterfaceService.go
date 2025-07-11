package service

import "github.com/charmbracelet/bubbles/list"

type BastionServiceInterface interface {
	List() ([]list.Item, error)
}
