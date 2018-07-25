package content

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func initTestServiceConfig(contentStoreURL string, contentPreviewURL string) ServiceConfig {
	return ServiceConfig{
		ContentStoreAppName:        "content-source-app",
		ContentStoreAppHealthURI:   contentStoreURL,
		ContentPreviewAppName:      "content-preview-app",
		ContentPreviewAppHealthURI: contentPreviewURL,
		HTTPClient:                 http.DefaultClient,
	}
}

func TestServiceConfig_ContentStoreCheck(t *testing.T) {
	ts := startFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig(ts.URL, "")

	check := sc.ContentStoreCheck()
	out, _ := check.Checker()
	assert.Equal(t, out, "Ok")
}

func TestServiceConfig_ContentStoreCheck_NotHealthy(t *testing.T) {
	ts := startNotFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig(ts.URL, "")

	check := sc.ContentStoreCheck()
	_, err := check.Checker()
	assert.Error(t, err, sc.ContentStoreAppName+" service is not responding with OK. Status=502")
}

func TestServiceConfig_ContentStoreCheck_InvalidAddress(t *testing.T) {
	sc := initTestServiceConfig("http://sampleHost:8080", "")
	check := sc.ContentStoreCheck()
	_, err := check.Checker()
	assert.Error(t, err, "dial tcp: lookup sampleHost: no such host")
}

func TestServiceConfig_ContentPreviewCheck(t *testing.T) {
	ts := startFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig("", ts.URL)

	check := sc.ContentPreviewCheck()
	out, _ := check.Checker()
	assert.Equal(t, out, "Ok")
}

func TestServiceConfig_ContentPreviewCheck_NotHealthy(t *testing.T) {
	ts := startNotFunctionalService()
	defer ts.Close()
	sc := initTestServiceConfig("", ts.URL)

	check := sc.ContentPreviewCheck()
	_, err := check.Checker()
	assert.Error(t, err, sc.ContentStoreAppName+" service is not responding with OK. Status=502")
}

func TestServiceConfig_ContentPreviewCheck_InvalidAddress(t *testing.T) {
	sc := initTestServiceConfig("http://sampleHost:8080", "")
	check := sc.ContentPreviewCheck()
	_, err := check.Checker()
	assert.Error(t, err, "dial tcp: lookup sampleHost: no such host")
}

func TestServiceConfig_GtgCheck(t *testing.T) {
	contentStoreTestService := startFunctionalService()
	defer contentStoreTestService.Close()

	contentPreviewTestService := startFunctionalService()
	defer contentPreviewTestService.Close()

	internalContentPreviewTestService := startFunctionalService()
	defer internalContentPreviewTestService.Close()

	sc := initTestServiceConfig(contentStoreTestService.URL, contentPreviewTestService.URL)

	status := sc.GtgCheck()
	assert.Equal(t, true, status.GoodToGo)
}

func TestServiceConfig_GtgCheck_NotGtg(t *testing.T) {
	contentStoreTestService := startNotFunctionalService()
	defer contentStoreTestService.Close()

	contentPreviewTestService := startNotFunctionalService()
	defer contentPreviewTestService.Close()

	internalContentPreviewTestService := startNotFunctionalService()
	defer internalContentPreviewTestService.Close()

	sc := initTestServiceConfig(contentStoreTestService.URL, contentPreviewTestService.URL)

	status := sc.GtgCheck()
	assert.Equal(t, false, status.GoodToGo)
}
