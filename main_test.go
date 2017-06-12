package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Financial-Times/image-resolver/content"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const contentSourceAppMame = "content-source-app-name"

var (
	imageResolver  *httptest.Server
	contentAPIMock *httptest.Server
)

func startContentAPIMock(contentApiMock func(http.ResponseWriter, *http.Request), healthMock func(http.ResponseWriter, *http.Request)) {
	router := mux.NewRouter()
	router.Path("/").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(contentApiMock)})
	router.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(healthMock)})
	router.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(healthMock)})
	contentAPIMock = httptest.NewServer(router)
}

func functionalEnrichedContentAPIMock(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("test-resources/valid-content-source-response.json")
	if err != nil {
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func functionalAPIMockForInternalContent(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("test-resources/valid-internalcontent-source-response.json")
	if err != nil {
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func contentApiStatusErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func contentApiStatusOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func startImageResolverService() {
	sc := content.ServiceConfig{
		ContentSourceAppName: contentSourceAppMame,
		ContentSourceURL:     contentAPIMock.URL,
		HttpClient:           http.DefaultClient,
	}

	r := content.NewContentReader("content-source-app-name", contentAPIMock.URL, http.DefaultClient)
	ir := content.NewImageResolver(r, "http://www.ft.com/ontology/content/ImageSet", "test.api.ft.com")

	h := setupServiceHandler(ir, sc)
	imageResolver = httptest.NewServer(h)
}

func stopServices() {
	imageResolver.Close()
	contentAPIMock.Close()
}

func TestShouldReturn200ForContentImages(t *testing.T) {
	startContentAPIMock(functionalEnrichedContentAPIMock, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	expected, err := ioutil.ReadFile("test-resources/valid-expanded-content-response.json")
	assert.NoError(t, err, "")

	body, err := ioutil.ReadFile("test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read file necessary for test case")
	resp, err := http.Post(imageResolver.URL+"/content/image", "application/json", bytes.NewReader(body))
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	actualResponse, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "")

	assert.JSONEq(t, string(expected), string(actualResponse))
}

func TestShouldReturn200ForInternalContentImages(t *testing.T) {
	startContentAPIMock(functionalAPIMockForInternalContent, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	expected, err := ioutil.ReadFile("test-resources/valid-expanded-internalcontent-response.json")
	assert.NoError(t, err, "")

	body, err := ioutil.ReadFile("test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read file necessary for test case")
	resp, err := http.Post(imageResolver.URL+"/internalcontent/image", "application/json", bytes.NewReader(body))
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	actualResponse, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "")

	assert.JSONEq(t, string(expected), string(actualResponse))
}

func TestShouldReturn400InvalidJsonContentEndpoint(t *testing.T) {
	startContentAPIMock(functionalEnrichedContentAPIMock, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	body := `{"body":"invalid""body"}`
	resp, err := http.Post(imageResolver.URL+"/content/image", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShouldReturn500InvalidJsonContentEndpoint(t *testing.T) {
	startContentAPIMock(functionalEnrichedContentAPIMock, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	body := `{"id":"36037ab1-da3b-35bf-b5ee-4fc23723b635"}`
	resp, err := http.Post(imageResolver.URL+"/content/image", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestShouldReturn400InvalidJsonInternalContentEndpoint(t *testing.T) {
	startContentAPIMock(functionalAPIMockForInternalContent, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	body := `{"body":"invalid""body"}`
	resp, err := http.Post(imageResolver.URL+"/internalcontent/image", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShouldReturn500InvalidJsonInternalContentEndpoint(t *testing.T) {
	startContentAPIMock(functionalEnrichedContentAPIMock, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	body := `{"id":"36037ab1-da3b-35bf-b5ee-4fc23723b635"}`
	resp, err := http.Post(imageResolver.URL+"/internalcontent/image", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestShouldBeHealthy(t *testing.T) {
	startContentAPIMock(functionalEnrichedContentAPIMock, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__health")
	assert.NoError(t, err, "Cannot send request to health endpoint")

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeHealthyWhenContentApiIsNotHappy(t *testing.T) {
	startContentAPIMock(contentApiStatusErrorHandler, contentApiStatusErrorHandler)
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__health")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 503")
	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "")
	assertMsg := fmt.Sprintf(`"id":"check-connect-%v","name":"Check connectivity to %v","ok":false`, contentSourceAppMame, contentSourceAppMame)
	assert.Contains(t, string(respBody), assertMsg)

}

func TestShouldBeGoodToGo(t *testing.T) {
	startContentAPIMock(functionalEnrichedContentAPIMock, contentApiStatusOkHandler)
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeGoodToGoWhenContentApiIsNotHappy(t *testing.T) {
	startContentAPIMock(contentApiStatusErrorHandler, contentApiStatusErrorHandler)
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}
