package data

import (
	"errors"
	"fmt"
)

type Item struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type Items []Item

var items Items

func init() {
	items = Items{
		Item{
			Id:    0,
			Name:  "Table",
			Price: 100.00,
		},
		Item{
			Id:    1,
			Name:  "Chair",
			Price: 49.99,
		},
	}
}

func GetItems() Items {
	return items
}

func GetItem(id int) (Item, error) {
	for _, item := range items {
		if item.Id == id {
			return item, nil
		}
	}
	return Item{
		Id:    -1,
		Name:  "",
		Price: -1.00,
	}, errors.New(fmt.Sprintf("item does not exist with id=%d", id))
}

func AddItem(name string, price float32) (err error) {
	if name == "" || price <= 0.0 {
		return errors.New(
			fmt.Sprintf("invalid item data: %v", []any{name, price}))
	}
	newItem := Item{
		Id:    items[len(items)-1].Id + 1,
		Name:  name,
		Price: price,
	}
	items = append(items, newItem)
	return nil
}
