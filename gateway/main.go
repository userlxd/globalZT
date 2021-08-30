package main

import (
	"globalZT/tools/log"
)

func main() {
	log.Log.Info("start")

	Run()

	defer log.Log.Info("quit")
}
