package main

import (
  "fmt"
  "runtime"
  "github.com/frostwind/l2go/loginserver"
)

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  loginserver.Init()

  fmt.Println("Server stopped.")
}
