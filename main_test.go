package main

import (
	"bytes"
	"encoding/json"
	"github.com/Financial-Times/image-resolver/content"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var imageResolver *httptest.Server
var contentAPIMock *httptest.Server

func startContentAPIMock(status string) {
	router := mux.NewRouter()
	var getContent http.HandlerFunc
	var health http.HandlerFunc

	if status == "happy" {
		getContent = happyContentAPIMock
		health = happyHandler
	} else if status == "badRequest" {
		getContent = badRequest
		health = happyHandler
	} else {
		getContent = internalErrorHandler
		health = internalErrorHandler
	}

	router.Path("/content/{uuid}").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(getContent)})
	router.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(health)})
	router.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(health)})

	contentAPIMock = httptest.NewServer(router)
}

func happyContentAPIMock(writer http.ResponseWriter, request *http.Request) {
	file, err := os.Open("test-resources/content.json")
	if err != nil {
		return
	}
	defer file.Close()
	io.Copy(writer, file)
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func happyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func startImageResolverService() {
	contentAPIURI := contentAPIMock.URL + "/content/"
	router := strings.Replace(contentAPIMock.URL, "http://", "", -1)
	sc := content.ServiceConfig{
		"9090",
		"content-public-read",
		router,
		"",
		"",
		http.DefaultClient,
	}

	var reader content.Reader
	var parser content.Parser
	var ir content.ImageResolver

	reader = content.NewContentReader(contentAPIURI, router)
	parser = content.NewBodyParser(content.ImageSetType)
	ir = *content.NewImageResolver(&reader, &parser)
	appLogger := content.NewAppLogger()
	contentHandler := content.ContentHandler{&sc, &ir, appLogger}
	h := setupServiceHandler(sc, &contentHandler)
	imageResolver = httptest.NewServer(h)
}

func stopServices() {
	imageResolver.Close()
	contentAPIMock.Close()
}

func getMapFromReader(r io.Reader) map[string]interface{} {
	var m map[string]interface{}
	json.NewDecoder(r).Decode(&m)
	return m
}

func TestShouldReturn200Content(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()
	file, _ := os.Open("test-resources/content.json")
	defer file.Close()
	resp, err := http.Get(contentAPIMock.URL + "/content/22c0d426-1466-11e7-b0c1-37e417ee6c76")
	assert.NoError(t, err, "Cannot send request to content endpoint")
	defer resp.Body.Close()

	expectedOutput := getMapFromReader(file)
	actualOutput := getMapFromReader(resp.Body)

	assert.Equal(t, expectedOutput, actualOutput, "Response body shoud be equal to transformer response body")
	var jsonStr = []byte(`{"id":"22c0d426-1466-11e7-b0c1-37e417ee6c76"}`)
	respPost, errPost := http.Post(imageResolver.URL+"/content/image", "application/json", bytes.NewBuffer(jsonStr))
	assert.NoError(t, errPost, "Cannot send request to imageresolver /content endpoint")
	defer respPost.Body.Close()

	assert.Equal(t, http.StatusOK, respPost.StatusCode, "Response status should be 200")
}

func TestShouldReturn200LeadImages(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()
	file, _ := os.Open("test-resources/content.json")
	defer file.Close()
	resp, err := http.Get(contentAPIMock.URL + "/content/22c0d426-1466-11e7-b0c1-37e417ee6c76")
	assert.NoError(t, err, "Cannot send request to content endpoint")
	defer resp.Body.Close()

	expectedOutput := getMapFromReader(file)
	actualOutput := getMapFromReader(resp.Body)

	assert.Equal(t, expectedOutput, actualOutput, "Response body shoud be equal to transformer response body")
	var jsonStr = []byte(`{"uuid":"22c0d426-1466-11e7-b0c1-37e417ee6c76"}`)
	respPost, errPost := http.Post(imageResolver.URL+"/internalcontent/image", "application/json", bytes.NewBuffer(jsonStr))
	assert.NoError(t, errPost, "Cannot send request to imageresolver /internalcontent endpoint")
	defer respPost.Body.Close()

	assert.Equal(t, http.StatusOK, respPost.StatusCode, "Response status should be 200")
}

func TestShouldReturn400InvalidJsonLeadImages(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()
	var invalidjsonStr = []byte(`{
	  "design": null,
	  "tableOfContents": null,
	  "topper": null,
	  "leadImages": [
	    {
	      "id": "89f194c8-13bc-11e7-80f4-13e067d5072c",
	      "type": "square"
	    },
	    {
	      "id": "3e96c818-13bc-11e7-b0c1-37e417ee6c76",
	      "type": "standard"
	    },
	    {
	      "id": "8d7b4e22-13bc-11e7-80f4-13e067d5072c",
	      "type": "wide"
	    }
	  ],
	  "uuid": "5010e2e4-09bd-11e7-97d1-5e720a26771b",
	  "lastModified": "2017-03-31T08:23:37.061Z",
	  "publishReference":
	}`)
	resp, err := http.Post(imageResolver.URL+"/internalcontent/image", "", bytes.NewBuffer(invalidjsonStr))
	assert.NoError(t, err, "Cannot send request to content endpoint")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response status should be 400")
}

func TestShouldReturn400InvalidJsonContent(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()
	var invalidjsonStr = []byte(`{
	  "id": "http://www.ft.com/thing/639cd952-149f-11e7-2ea7-a07ecd9ac73f",
	  "type": "http://www.ft.com/ontology/content/ImageSet",
	  "title": "",
	  "alternativeTitles": {},
	  "alternativeStandfirsts": {},
	  "publishReference": "tid_5ypvntzcpu",
	  "lastModified": "2017-03-29T19:39:18.226Z",
	  "canBeDistributed":
	  }`)
	resp, err := http.Post(imageResolver.URL+"/content/image", "", bytes.NewBuffer(invalidjsonStr))
	assert.NoError(t, err, "Cannot send request to content endpoint")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response status should be 400")
}

func TestShouldBeHealthy(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__health")
	assert.NoError(t, err, "Cannot send request to health endpoint")

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeHealthyWhenContentApiIsNotHappy(t *testing.T) {
	startContentAPIMock("unhappy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(contentAPIMock.URL + "/__health")
	assert.NoError(t, err, "Cannot send request to health endpoint")
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}

func TestShouldBeGoodToGo(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeGoodToGoWhenContentApiIsNotHappy(t *testing.T) {
	startContentAPIMock("unhappy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(contentAPIMock.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}
