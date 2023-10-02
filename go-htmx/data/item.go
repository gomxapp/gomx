package data

import (
	"errors"
	"fmt"
)

type Item struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type Items []Item

var items Items = make(Items, 0)

func SeedData() {
	items = Items{
		Item{
			Name:  "Table",
			Price: 100.00,
		},
		Item{
			Name:  "Chair",
			Price: 49.99,
		},
	}
}

func GetItems() Items {
	return items
}

func AddItem(newItem Item) (count int, err error) {
	if newItem.Name == "" || newItem.Price <= 0.0 {
		return 0, errors.New(fmt.Sprintf("Invalid item data: %v", newItem))
	}
	items = append(items, newItem)
	return len(items), nil
}
