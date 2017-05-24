package content

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func startFunctionalService() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func startNotFunctionalService() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
}

func initTestServiceConfig(URL string) ServiceConfig {
	return ServiceConfig{
		ContentSourceAppName: "content-source-app",
		ContentSourceURL:     URL,
		HttpClient:           http.DefaultClient,
	}
}

func TestServiceConfig_ContentCheck(t *testing.T) {
	ts := startFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig(ts.URL)

	check := sc.ContentCheck()
	out, _ := check.Checker()
	assert.Equal(t, out, "Ok")
}

func TestServiceConfig_ContentCheck_NotHealthy(t *testing.T) {
	ts := startNotFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig(ts.URL)

	check := sc.ContentCheck()
	_, err := check.Checker()
	assert.Error(t, err, sc.ContentSourceAppName+" service is not responding with OK. Status=502")
}

func TestServiceConfig_ContentCheck_InvalidAddress(t *testing.T) {
	sc := initTestServiceConfig("http://sampleHost:8080")
	check := sc.ContentCheck()
	_, err := check.Checker()
	assert.Error(t, err, "dial tcp: lookup sampleHost: no such host")
}

func TestServiceConfig_GtgCheck(t *testing.T) {
	ts := startFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig(ts.URL)

	status := sc.GtgCheck()
	assert.Equal(t, status.GoodToGo, true)
}

func TestServiceConfig_GtgCheck_NotGtg(t *testing.T) {
	ts := startNotFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig(ts.URL)

	status := sc.GtgCheck()
	assert.Equal(t, status.GoodToGo, false)
}
