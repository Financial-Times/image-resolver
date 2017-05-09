package content

import (
	"errors"
	"fmt"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
	"net/http"
)

func (sc *ServiceConfig) GtgCheck() gtg.Status {
	msg, err := sc.checkerContent()
	if err != nil {
		return gtg.Status{GoodToGo: false, Message: msg}
	}

	return gtg.Status{GoodToGo: true}
}

func (sc *ServiceConfig) ContentCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "check-connect-content-public-read",
		Name:             "Check connectivity to content-public-read",
		Severity:         1,
		BusinessImpact:   "Image unrolled won't be available",
		TechnicalSummary: fmt.Sprintf(`Cannot connect to content-public-read.`),
		PanicGuide:       "https://dewey.ft.com/upp-image-resolver.html",
		Checker: func() (string, error) {
			return sc.checkerContent()
		},
	}
}

func (sc *ServiceConfig) checkerContent() (string, error) {
	healthUri := "http://" + sc.RouterAddress + "/__health"
	req, err := http.NewRequest("GET", healthUri, nil)
	req.Host = sc.Content_public_read
	resp, err := sc.HttpClient.Do(req)
	if err != nil {
		msg := fmt.Sprintf("%s service is unreachable: %v", "content-public-read", err)
		return msg, errors.New(msg)
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("%s service is not responding with OK. status=%d", "content-public-read", resp.StatusCode)
		return msg, errors.New(msg)
	}
	return "Ok", nil
}
