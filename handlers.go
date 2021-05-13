package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"github.com/NOVAPokemon/utils/clients"
	"github.com/NOVAPokemon/utils/items"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/NOVAPokemon/utils/websockets"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const itemsFile = "store_items.json"

var (
	itemsMap       map[string]items.StoreItem
	marshaledItems []byte

	httpClient = &http.Client{
		Timeout:   websockets.Timeout,
		Transport: clients.NewTransport(),
	}
	basicClient = clients.NewBasicClient(false, "")

	serverName   string
	commsManager websockets.CommunicationManager
)

func init() {
	if aux, exists := os.LookupEnv(utils.HostnameEnvVar); exists {
		serverName = aux
	} else {
		log.Fatal("could not load server name")
	}
	log.Info("Server name : ", serverName)

	var err error
	itemsMap, marshaledItems, err = loadShopItems()
	if err != nil {
		log.Fatal(err)
	}
}

func handleGetItems(w http.ResponseWriter, r *http.Request) {
	_, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetItemsError(err), http.StatusBadRequest)
		return
	}

	_, err = w.Write(marshaledItems)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetItemsError(err), http.StatusInternalServerError)
	}
}

func handleBuyItem(w http.ResponseWriter, r *http.Request) {
	itemName := mux.Vars(r)[api.ShopItemNameVar]

	toBuy, ok := itemsMap[itemName]
	if !ok {
		err := wrapBuyItemError(newItemNotFoundError(itemName))
		utils.LogAndSendHTTPError(&w, err, http.StatusNotFound)
		return
	}

	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapBuyItemError(err), http.StatusUnauthorized)
		return
	}

	trainerStatsToken, err := tokens.ExtractAndVerifyTrainerStatsToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapBuyItemError(err), http.StatusUnauthorized)
		return
	}

	authTokenString := r.Header.Get(tokens.AuthTokenHeaderName)
	if trainerStatsToken.TrainerStats.Coins < toBuy.Price {
		err = wrapBuyItemError(warnNotEnoughMoney)
		utils.LogWarnAndSendHTTPError(&w, err, http.StatusForbidden)
		return
	}

	item := toBuy.ToItem()
	toAdd := []items.Item{item}

	trainersClient := clients.NewTrainersClient(httpClient, commsManager, basicClient)

	_, err = trainersClient.AddItems(authToken.Username, toAdd, authTokenString)
	if err != nil {
		utils.LogAndSendHTTPError(&w, err, http.StatusInternalServerError)
		return
	}

	newTrainerStats := trainerStatsToken.TrainerStats
	newTrainerStats.Coins -= toBuy.Price
	_, err = trainersClient.UpdateTrainerStats(authToken.Username, newTrainerStats, authTokenString)
	if err != nil {
		utils.LogAndSendHTTPError(&w, err, http.StatusInternalServerError)
		return
	}

	log.Info("trainer items and money updated")

	w.Header().Set(tokens.ItemsTokenHeaderName, trainersClient.ItemsToken)
	w.Header().Set(tokens.StatsTokenHeaderName, trainersClient.TrainerStatsToken)
}

func loadShopItems() (map[string]items.StoreItem, []byte, error) {
	data, err := ioutil.ReadFile(itemsFile)
	if err != nil {
		return nil, nil, wrapLoadShopItemsError(err)
	}

	var itemsArr []items.StoreItem
	err = json.Unmarshal(data, &itemsArr)
	if err != nil {
		return nil, nil, wrapLoadShopItemsError(err)
	}

	itemsMapAux := make(map[string]items.StoreItem, len(itemsArr))
	for _, item := range itemsArr {
		itemsMapAux[item.Name] = item
	}

	log.Infof("Loaded %d items.", len(itemsArr))

	marshalledItems, err := json.Marshal(itemsArr)
	if err != nil {
		return nil, nil, wrapLoadShopItemsError(err)
	}

	return itemsMapAux, marshalledItems, nil
}
