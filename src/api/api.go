/*
Package api for the overlord software.
This contains functions that use the atreugo
golang library to create the HTTP REST API
that may be enabled when running overlord.
*/
package api

import (
	"bitsnthings.dev/overlord/src/api/middlewarez"
	"bitsnthings.dev/overlord/src/api/routes"
	"bitsnthings.dev/overlord/src/log"
	"github.com/savsgio/atreugo/v11"
)

type API struct {
	server *atreugo.Atreugo
	Conf   *atreugo.Config
}

func (api *API) Setup(conf *atreugo.Config) {
	if conf != nil {
		api.Conf = conf
	}
	api.server = atreugo.New(*api.Conf)
	middlewarez.SetMiddlewares(api.server)
	routes.SetRoutes(api.server)
}

func (api *API) Start() {
	if err := api.server.ListenAndServe(); err != nil {
		log.PrintLog(
			log.ERROR,
			"Error starting atreugo HTTP API! %s",
			err,
		)
	}
}
