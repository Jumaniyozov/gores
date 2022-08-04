package api

import (
	"github.com/gorilla/mux"
	"github.com/jumaniyozov/gores/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

type API struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *Config) *API {
	return &API{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (api *API) Start() error {
	if err := api.configerLoggerField(); err != nil {
		return err
	}

	api.logger.Info("Starting api server at port", api.config.BindAddr)

	api.configerRouterField()
	if err := api.configerStorageField(); err != nil {
		return err
	}

	return http.ListenAndServe(api.config.BindAddr, api.router)
}
