package main

import (
	"errorwrapper/errorspec"
	"errorwrapper/wrapper"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.WithField("svc-name", "error-wrapper-test")
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	a := wrapper.New(errorspec.DatabaseFailure).Errorln("root level error")
	b := wrapper.New(a).WithFields(wrapper.Fields{
		"resource-id":    "some additional metadata",
		"database-table": "oauth.Customer",
	}).Errorln("middle layer error")
	c := wrapper.New(b).WithField("endpoint", "/resource").Errorln("top level error")

	logger.WithError(c).WithField("trace", wrapper.Unfold(c)).Errorln("Everything is bad")
}
