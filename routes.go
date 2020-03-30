package main

import (
	"github.com/NOVAPokemon/utils"
)


const GET = "GET"
const POST = "POST"

const GetShopItems = "GET_ITEMS"
const GetShopItemsPath = "/shop/items/"

const BuyItem = "BUY_ITEM"
const BuyItemPath = "/shop/items/buy/{itemId}"

var routes = utils.Routes{

	utils.Route{
		Name:        GetShopItems,
		Method:      GET,
		Pattern:     GetShopItemsPath,
		HandlerFunc: HandleGetItems,
	},

	utils.Route{
		Name:        BuyItem,
		Method:      GET,
		Pattern:     BuyItemPath,
		HandlerFunc: HandleBuyItem,
	},
}
