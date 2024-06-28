package gosyga

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type apiWithLogger struct {
	// for sync gateway authorization
	user     string
	password string

	log *logrus.Entry
}

func newNullApiLogger(user, password string) apiWithLogger {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	return apiWithLogger{
		user:     user,
		password: password,

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
