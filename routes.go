package main

import (
	"fmt"
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const (
	getItemsName = "GET_ITEMS"
	buyItemName  = "BUY_ITEM"
)

const (
	get  = "GET"
	post = "POST"
)

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(fmt.Sprintf("%s", serviceName))),
	utils.Route{
		Name:        getItemsName,
		Method:      get,
		Pattern:     api.GetShopItemsPath,
		HandlerFunc: handleGetItems,
	},

	utils.Route{
		Name:        buyItemName,
		Method:      post,
		Pattern:     api.BuyItemsRoute,
		HandlerFunc: handleBuyItem,
	},
}
