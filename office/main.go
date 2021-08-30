package main

import (
	"globalZT/tools/log"
)

func main() {
	log.Log.Info("start")
	defer log.Log.Info("quit")
}
