package main

import (
	"encoding/json"
	"errors"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/cookies"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const ItemsFile = "store_items.json"

var itemsArr, itemsMap = loadShopItems()
var marshaledItems, _ = json.Marshal(itemsArr)

var ErrItemNotFound = errors.New("Item Not Found")
var ErrTrainerStatsTokenNotFound = errors.New("Trainer stats token not found")
var ErrNotEnoughMoney = errors.New("Item Not Found")

func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write(marshaledItems)
}

func HandleBuyItem(w http.ResponseWriter, r *http.Request) {

	itemName := mux.Vars(r)["itemName"]

	toBuy, ok := itemsMap[itemName]

	if !ok {
		http.Error(w, ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}

	_, err := cookies.ExtractAndVerifyAuthToken(&w, r, "store")
	if err != nil {
		return
	}

	trainerStatsToken, err := cookies.ExtractTrainerStatsToken(r)

	if err != nil {
		http.Error(w, ErrTrainerStatsTokenNotFound.Error(), http.StatusUnauthorized)
	}

	if trainerStatsToken.TrainerStats.Coins >= toBuy.Price {
		//TODO commit buy to trainers
	} else {
		http.Error(w, ErrNotEnoughMoney.Error(), http.StatusForbidden)
	}

}

func loadShopItems() ([]utils.StoreItem, map[string]utils.StoreItem) {

	data, err := ioutil.ReadFile(ItemsFile)
	if err != nil {
		log.Errorf("Error loading items file ")
		log.Fatal(err)
		panic(err)
	}

	var itemsArr []utils.StoreItem
	err = json.Unmarshal(data, &itemsArr)

	var itemsMap = make(map[string]utils.StoreItem, len(itemsArr))
	for _, item := range itemsArr {
		itemsMap[item.Name] = item
	}

	if err != nil {
		log.Errorf("Error unmarshalling item names")
		log.Fatal(err)
		panic(err)
	}

	log.Infof("Loaded %d items.", len(itemsArr))

	return itemsArr, itemsMap
}
