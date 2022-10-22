package main

import (
	"go-kubernetes/http"
	"go-kubernetes/logger"
	"sync"
)

func main() {
	// Initialize the logger
	Logger := logger.Logger
	defer Logger.Sync()
	Logger.Info("Application starting")

	http.InitServer()
	http.InitHealth()
	http.InitMonitors()
	http.InitControllers()
	http.InitWebsocket()
	http.InitWebsocketControllers()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go http.Start(":3000", &wg)

	wg.Wait()
}
