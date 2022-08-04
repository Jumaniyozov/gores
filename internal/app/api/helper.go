package api

import (
	"github.com/jumaniyozov/gores/storage"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func (api *API) configerLoggerField() error {
	logLevel, err := logrus.ParseLevel(
		api.config.LoggerLevel,
	)
	if err != nil {
		return err
	}

	api.logger.SetLevel(logLevel)
	return nil
}

func (api *API) configerRouterField() {
	api.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello! RestApi"))
		if err != nil {
			log.Fatal(err)
		}
	})
}

func (api *API) configerStorageField() error {
	storageDB := storage.New(api.config.Storage)

	if err := storageDB.Open(); err != nil {
		return err
	}

	api.storage = storageDB
	api.logger.Info("Database connection successfully created!")

	return nil
}
