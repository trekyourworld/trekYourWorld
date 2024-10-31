package server

import (
	"trekyourworld/handlers"

	"oss.nandlabs.io/golly/lifecycle"
	"oss.nandlabs.io/golly/rest/server"
)

func Start() {
	srvr, err := server.Default()
	if err != nil {
		panic(err)
	}

	srvr.Opts().PathPrefix = "/api/v1"

	srvr.Get("/orgs", handlers.FindAllOrganisations)

	srvr.Get("/treks", handlers.FindAllTreks)
	srvr.Get("/treks/search", handlers.SearchTrek)
	srvr.Post("/treks/filter", handlers.FilterTreks)

	srvr.Post("/contact", handlers.ContactUs)

	manager := lifecycle.NewSimpleComponentManager()
	manager.Register(srvr)

	manager.StartAndWait()
}
