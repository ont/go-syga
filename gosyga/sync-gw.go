package gosyga

import "github.com/sirupsen/logrus"

type apiWithLogger struct {
	logger *logrus.Entry
}

func (a *AdminApi) WithLogger(logger *logrus.Entry) *AdminApi {
	a.logger = logger
	return a
}

func (c *ClientApi) WithLogger(logger *logrus.Entry) *ClientApi {
	c.logger = logger
	return c
}
