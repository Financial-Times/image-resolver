package content

import (
	"github.com/Sirupsen/logrus"
	"fmt"
)

type appLogger struct {
	log *logrus.Logger
}

func NewAppLogger() *appLogger {
	logrus.SetLevel(logrus.InfoLevel)
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	return &appLogger{log: log}
}

func (appLogger *appLogger) TransactionStartedEvent(requestURL string, transactionID string, uuid string) {
	appLogger.log.WithFields(logrus.Fields{
		"request_url":    requestURL,
		"transaction_id": transactionID,
		"uuid":           uuid,
	}).Infof("Transaction started %s", transactionID)
}

func (appLogger *appLogger) TransactionFinishedEvent(requestURL string, transactionID string, statusCode int, uuid string, message string) {
	e := appLogger.log.WithFields(logrus.Fields{
		"request_url":    requestURL,
		"transaction_id": transactionID,
		"uuid":           uuid,
	})

	msg := fmt.Sprintf("Transaction %s finished with status %d: %s", transactionID, statusCode, message)
	if statusCode < 300 {
		e.Infof(msg)
	} else {
		e.Errorf(msg)
	}
}

func (appLogger *appLogger) Infof(uuid string, format string, args ...interface{}) {
	appLogger.log.WithFields(logrus.Fields{"uuid": uuid}).Infof(format, args)
}
