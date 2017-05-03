package main
import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/Financial-Times/image-resolver/content"
	"strings"
	"bytes"
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
	router := strings.Replace(contentAPIMock.URL, "http://", "",-1)
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
	parser = content.BodyParser{}
	ir = *content.NewImageResolver(&reader, &parser)

	contentHandler := content.ContentHandler{&sc, &ir}
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

func TestShouldReturn200(t *testing.T) {
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
	respPost, errPost := http.Post(imageResolver.URL + "/content", "application/json", bytes.NewBuffer(jsonStr))
	assert.NoError(t, errPost, "Cannot send request to imageresolver endpoint")
	defer respPost.Body.Close()

	assert.Equal(t, http.StatusOK, respPost.StatusCode, "Response status should be 200")
}

func TestShouldReturn400InvalidJson(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()
	var jsonStr = []byte(`{
	"id": "http://www.ft.com/thing/22c0d426-1466-11e7-b0c1-37e417ee6c76",
		"type": "http://www.ft.com/ontology/content/Article",
		"blabla"}`)
	resp, err := http.Post(imageResolver.URL + "/content", "",bytes.NewBuffer(jsonStr))
	assert.NoError(t, err, "Cannot send request to content endpoint")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response status should be 400")
}

func TestShouldReturn400InvalidID(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()
	var jsonStr = []byte(`{"id":"22c0d426-1466-11e7-b0c1-37e417ee6c76xxxxx"}`)
	respPost, errPost := http.Post(imageResolver.URL + "/content", "application/json", bytes.NewBuffer(jsonStr))
	assert.NoError(t, errPost, "Cannot send request to content endpoint")
	defer respPost.Body.Close()

	assert.Equal(t, http.StatusBadRequest, respPost.StatusCode, "Response status should be 400")
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