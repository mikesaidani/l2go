package main

import (
	"fmt"
	"github.com/frostwind/l2go/loginserver"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	loginserver.Init()

	fmt.Println("Server stopped.")
}
