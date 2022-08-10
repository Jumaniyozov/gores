package api

import (
	"github.com/jumaniyozov/gores/storage"
	"github.com/sirupsen/logrus"
)

var (
	prefix string = "/api/v1"
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
	api.router.HandleFunc(prefix+"/articles", api.GetAllArticles).Methods("GET")
	api.router.HandleFunc(prefix+"/articles/{id}", api.GetArticleByID).Methods("GET")
	api.router.HandleFunc(prefix+"/articles/{id}", api.DeleteArticleByID).Methods("DELETE")
	api.router.HandleFunc(prefix+"/articles", api.PostArticle).Methods("POST")
	api.router.HandleFunc(prefix+"/user/register", api.PostUserRegister).Methods("POST")
}

func (api *API) configerStorageField() error {
	storageDB := storage.New(api.config.Storage)

	if err := storageDB.Open(); err != nil {
		return err
	}

	api.storage = storageDB

	return nil
}
