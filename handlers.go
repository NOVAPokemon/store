package main

import (
	"encoding/json"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const ItemsFile = "store_items.json"
var items = loadShopItems()
var marshaledItems, _ = json.Marshal(items)

func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write(marshaledItems)
}

func HandleBuyItem(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func loadShopItems() []utils.StoreItem {

	data, err := ioutil.ReadFile(ItemsFile)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var itemNames []utils.StoreItem
	err = json.Unmarshal(data, &itemNames)

	if err != nil {
		log.Errorf("Error unmarshalling item names")
		log.Fatal(err)
	}

	log.Infof("Loaded %d items.", len(itemNames))
	return itemNames
}
