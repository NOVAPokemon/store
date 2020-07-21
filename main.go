package main

import (
	"github.com/NOVAPokemon/utils"
)

const (
	host        = utils.ServeHost
	port        = utils.StorePort
	serviceName = "STORE"
)

func main() {
	flags := utils.ParseFlags(serverName)

	if !*flags.LogToStdout {
		utils.SetLogFile(serverName)
	}

	if !*flags.DelayedComms {
		commsManager = utils.CreateDefaultCommunicationManager()
	} else {
		locationTag := utils.GetLocationTag(utils.DefaultLocationTagsFilename, serverName)
		commsManager = utils.CreateDelayedCommunicationManager(utils.DefaultDelayConfigFilename,
			locationTag, false)
	}

	utils.StartServer(serviceName, host, port, routes, commsManager)
}
