package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const host = utils.Host
const port = utils.StorePort

var addr = fmt.Sprintf("%s:%d", host, port)

func main() {
	rand.Seed(time.Now().Unix())
	r := utils.NewRouter(routes)
	log.Infof("Starting STORE server in port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
