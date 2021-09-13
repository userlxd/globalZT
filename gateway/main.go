package main

import (
	"globalZT/pkg/tunnel"
	"globalZT/tools/log"
)

func main() {
	log.Log.Info("start")
	tunnel.Run()
	defer log.Log.Info("quit")
}
