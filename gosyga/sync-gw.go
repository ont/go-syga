package gosyga

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type apiWithLogger struct {
	log *logrus.Entry
}

func newNullApiLogger() apiWithLogger {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	return apiWithLogger{
		log: logrus.NewEntry(logger),
	}
}

func (a *AdminApi) WithLogger(log *logrus.Entry) *AdminApi {
	a.log = log
	return a
}

func (c *ClientApi) WithLogger(log *logrus.Entry) *ClientApi {
	c.log = log
	return c
}
