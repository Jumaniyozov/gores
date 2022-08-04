package api

import "github.com/sirupsen/logrus"

func (api *API) configerLoggerField() error {
	logLevel, err := logrus.ParseLevel(api.config.LoggerLevel)
	if err != nil {
		return err
	}

	api.logger.SetLevel(logLevel)
	return nil
}
