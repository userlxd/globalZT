package main

import (
	"globalZT/office/office"
	"globalZT/tools/log"
)

func main() {
	log.Log.Info("start")
	office.Run()
	defer log.Log.Info("quit")
}
