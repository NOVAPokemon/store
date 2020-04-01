package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const GET = "GET"
const POST = "POST"

const GetShopItems = "GET_ITEMS"
const BuyItem = "BUY_ITEM"

var routes = utils.Routes{

	utils.Route{
		Name:        GetShopItems,
		Method:      GET,
		Pattern:     api.GetShopItemsPath,
		HandlerFunc: HandleGetItems,
	},

	utils.Route{
		Name:        BuyItem,
		Method:      GET,
		Pattern:     api.BuyItemPath,
		HandlerFunc: HandleBuyItem,
	},
}
