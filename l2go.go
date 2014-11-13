package main

import (
	"flag"
	"fmt"
	"github.com/frostwind/l2go/gameserver"
	"github.com/frostwind/l2go/loginserver"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var mode int
	flag.IntVar(&mode, "mode", 0, "Set to 0 to run the Login Server or 1 to run the Game Server")
	flag.Parse()

	if mode == 0 {
		loginserver.Init()
	} else {
		gameserver.Init()
	}

	fmt.Println("Server stopped.")
}
