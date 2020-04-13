package main

import (
	"github.com/NOVAPokemon/utils/items"
)

type StoreItem struct {
	Name  string
	Price int
}

func (storeItem StoreItem) ToItem() items.Item {
	return items.Item{
		Name:   storeItem.Name,
		Effect: items.GetEffectForItem(storeItem.Name),
	}
}
