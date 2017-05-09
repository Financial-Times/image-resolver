package content

import (
	"github.com/Sirupsen/logrus"
)

type event struct {
	serviceName   string
	requestURL    string
	transactionID string
	uuid          string
}

type appLogger struct {
	log *logrus.Logger
}

func NewAppLogger() *appLogger {
	logrus.SetLevel(logrus.InfoLevel)
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	return &appLogger{log}
}

func (appLogger *appLogger) ServiceStartedEvent(serviceName string, serviceConfig map[string]interface{}) {
	serviceConfig["event"] = "service_started"
	appLogger.log.WithFields(serviceConfig).Infof("%s started with configuration", serviceName)
}

func (appLogger *appLogger) TransactionStartedEvent(requestURL string, transactionID string, uuid string) {
	appLogger.log.WithFields(logrus.Fields{
		"event":          "transaction_started",
		"request_url":    requestURL,
		"transaction_id": transactionID,
		"uuid":           uuid,
	}).Info()
}

func (appLogger *appLogger) RequestEvent(serviceName string, requestURL string, transactionID string, uuid string) {
	appLogger.log.WithFields(logrus.Fields{
		"event":          "request",
		"service_name":   serviceName,
		"request_uri":    requestURL,
		"transaction_id": transactionID,
		"uuid":           uuid,
	}).Info()
}

func (appLogger *appLogger) RequestFailedEvent(serviceName string, requestURL string, transactionID string, statusCode int, uuid string) {
	appLogger.log.WithFields(logrus.Fields{
		"event":          "request_failed",
		"service_name":   serviceName,
		"request_url":    requestURL,
		"transaction_id": transactionID,
		"status":         statusCode,
		"uuid":           uuid,
	}).
		Warnf("Request failed. %s responded with %s", serviceName, statusCode)
}

func (appLogger *appLogger) ResponseEvent(serviceName string, requestURL string, transactionID string, statusCode int, uuid string) {
	appLogger.log.WithFields(logrus.Fields{
		"event":          "response",
		"service_name":   serviceName,
		"status":         statusCode,
		"request_url":    requestURL,
		"transaction_id": transactionID,
		"uuid":           uuid,
	}).
		Info("Response from " + serviceName)
}
