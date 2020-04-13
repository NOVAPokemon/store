package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"github.com/NOVAPokemon/utils/clients"
	"github.com/NOVAPokemon/utils/items"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const ItemsFile = "store_items.json"

var trainersClient = clients.NewTrainersClient(fmt.Sprintf("%s:%d", utils.Host, utils.TrainersPort))

var itemsArr, itemsMap = loadShopItems()
var marshaledItems, _ = json.Marshal(itemsArr)

var ErrItemNotFound = errors.New("item not Found")
var ErrNotEnoughMoney = errors.New("not enough money")
var ErrTrainerStatsTokenNotFound = errors.New("trainer stats token not found")
var ErrTrainerAuthTokenNotFound = errors.New("auth token not found")

func HandleGetItems(w http.ResponseWriter, r *http.Request) {
	_, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		return
	}

	_, _ = w.Write(marshaledItems)
}

func HandleBuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := mux.Vars(r)[api.ShopItemNameVar]

	toBuy, ok := itemsMap[itemName]
	if !ok {
		http.Error(w, ErrItemNotFound.Error(), http.StatusNotFound)
		return
	}

	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(ErrTrainerAuthTokenNotFound.Error())
		http.Error(w, ErrTrainerAuthTokenNotFound.Error(), http.StatusUnauthorized)
		return
	}

	trainerStatsToken, err := tokens.ExtractAndVerifyTrainerStatsToken(r.Header)
	if err != nil {
		log.Error(ErrTrainerStatsTokenNotFound.Error())
		http.Error(w, ErrTrainerStatsTokenNotFound.Error(), http.StatusUnauthorized)
		return
	}

	authTokenString := r.Header.Get(tokens.AuthTokenHeaderName)
	if trainerStatsToken.TrainerStats.Coins < toBuy.Price {
		log.Error(ErrNotEnoughMoney.Error())
		http.Error(w, ErrNotEnoughMoney.Error(), http.StatusForbidden)
		return
	}

	item := toBuy.ToItem()
	toAdd := []*items.Item{&item}
	added, err := trainersClient.AddItemsToBag(authToken.Username, toAdd, authTokenString)
	if err != nil {
		log.Error(err)
		return
	}

	if err := clients.CheckItemsAdded(toAdd, added); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Info("items were successfully added")
	}

	newTrainerStats := trainerStatsToken.TrainerStats
	newTrainerStats.Coins -= toBuy.Price
	trainerStats, err := trainersClient.UpdateTrainerStats(authToken.Username, newTrainerStats, authTokenString)
	if err != nil {
		log.Error(err)
		return
	}

	if err := clients.CheckUpdatedStats(&newTrainerStats, trainerStats); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		log.Info("stats were successfully updated")
	}

	err = trainersClient.GetItemsToken(authToken.Username, authTokenString)
	if err != nil {
		log.Error(err)
		return
	}

	w.Header().Set(tokens.ItemsTokenHeaderName, trainersClient.ItemsToken)
}

func loadShopItems() ([]items.StoreItem, map[string]items.StoreItem) {
	data, err := ioutil.ReadFile(ItemsFile)
	if err != nil {
		log.Errorf("Error loading items file ")
		log.Fatal(err)
		panic(err)
	}

	var itemsArr []items.StoreItem
	err = json.Unmarshal(data, &itemsArr)

	var itemsMap = make(map[string]items.StoreItem, len(itemsArr))
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
