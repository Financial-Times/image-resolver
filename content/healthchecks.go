package content

import (
	"fmt"
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
	"github.com/pkg/errors"
)

type ServiceConfig struct {
	ContentSourceAppName string
	ContentSourceURL     string
	HttpClient           *http.Client
}

func (sc *ServiceConfig) GtgCheck() gtg.Status {
	msg, err := sc.checkerContent()
	if err != nil {
		return gtg.Status{GoodToGo: false, Message: msg}
	}

	return gtg.Status{GoodToGo: true}
}

func (sc *ServiceConfig) ContentCheck() fthealth.Check {
	return fthealth.Check{
		ID:               fmt.Sprintf("check-connect-%s", sc.ContentSourceAppName),
		Name:             fmt.Sprintf("Check connectivity to %s", sc.ContentSourceAppName),
		Severity:         1,
		BusinessImpact:   "Image unrolled won't be available",
		TechnicalSummary: fmt.Sprintf(`Cannot connect to %v.`, sc.ContentSourceAppName),
		PanicGuide:       "https://dewey.ft.com/upp-image-resolver.html",
		Checker: func() (string, error) {
			return sc.checkerContent()
		},
	}
}

func (sc *ServiceConfig) checkerContent() (string, error) {
	healthUri := sc.ContentSourceURL + "/__health"
	req, err := http.NewRequest(http.MethodGet, healthUri, nil)
	req.Host = sc.ContentSourceAppName
	resp, err := sc.HttpClient.Do(req)
	if err != nil {
		return "Error", errors.Errorf("%s service is unreachable: %v", sc.ContentSourceAppName, err)
	}
	if resp.StatusCode != http.StatusOK {
		return "Error", errors.Errorf("%s service is not responding with OK. Status=%d", sc.ContentSourceAppName, resp.StatusCode)
	}
	return "Ok", nil
}
