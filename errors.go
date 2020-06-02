package main

import (
	"fmt"

	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	errorLoadShopItem = "error loading shop items"

	errorItemNotFoundFormat = "item %s not Found"
)

var (
	warnNotEnoughMoney = errors.New("not enough money")
)

// Handler wrappers
func wrapGetItemsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GetItemsName))
}

func wrapBuyItemError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, BuyItemName))
}

// Other wrappers
func wrapLoadShopItemsError(err error) error {
	return errors.Wrap(err, errorLoadShopItem)
}

// Error builders
func newItemNotFoundError(itemId string) error {
	return errors.New(fmt.Sprintf(errorItemNotFoundFormat, itemId))
}
