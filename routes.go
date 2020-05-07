package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const GET = "GET"
const POST = "POST"

const GetItemsName = "GET_ITEMS"
const BuyItemName = "BUY_ITEM"

var routes = utils.Routes{
	api.DefaultRoute,
	utils.Route{
		Name:        GetItemsName,
		Method:      GET,
		Pattern:     api.GetShopItemsPath,
		HandlerFunc: HandleGetItems,
	},

	utils.Route{
		Name:        BuyItemName,
		Method:      POST,
		Pattern:     api.BuyItemsRoute,
		HandlerFunc: HandleBuyItem,
	},
}
