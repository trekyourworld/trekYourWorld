package main

import (
	"trekyourworld/db"
	"trekyourworld/server"

	"oss.nandlabs.io/golly/l3"
)

var logger = l3.Get()

func main() {
	// init db
	err := db.Init()
	if err != nil {
		logger.ErrorF("Failed to initialize database: %v", err)
	} else {
		// start server
		server.Start()
	}
}
