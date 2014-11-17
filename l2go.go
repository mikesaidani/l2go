package main

import (
	"flag"
	"fmt"
	"github.com/frostwind/l2go/config"
	"github.com/frostwind/l2go/gameserver"
	"github.com/frostwind/l2go/loginserver"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var mode, gameServerId int
	flag.IntVar(&mode, "mode", 0, "Set to 0 to run the Login Server or 1 to run the Game Server")
	flag.IntVar(&gameServerId, "server", 1, "Set the id of the Game Server you want to run")
	flag.Parse()

	// Load the global configuration object
	globalConfig := config.Read()

	if mode == 0 {
		loginserver.Init(globalConfig)
	} else {
		// Load the Game Server specific configuration
		if gameServerId >= 1 && len(globalConfig.GameServers) >= gameServerId {
			var gameServerConfig config.GameServerConfigObject
			gameServerConfig.LoginServer = globalConfig.LoginServer
			gameServerConfig.GameServer = globalConfig.GameServers[gameServerId-1]
			gameserver.Init(gameServerConfig)
		} else {
			fmt.Println("No configuration found for the specified server.")
		}

	}

	fmt.Println("Server stopped.")
}
