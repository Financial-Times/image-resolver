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
	} else if status == "notFound" {
		getContent = notFoundHandler
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
	file, err := os.Open("test-resources/image-output.json")
	if err != nil {
		return
	}
	defer file.Close()
	io.Copy(writer, file)
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
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
	file, _ := os.Open("test-resources/image-output.json")
	defer file.Close()
	resp, err := http.Get(contentAPIMock.URL + "/content/5c3cae78-dbef-11e6-9d7c-be108f1c1dce")
	if err != nil {
		assert.FailNow(t, "Cannot send request to content endpoint", err.Error())
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")

	expectedOutput := getMapFromReader(file)
	actualOutput := getMapFromReader(resp.Body)

	assert.Equal(t, expectedOutput, actualOutput, "Response body shoud be equal to transformer response body")
}

func TestShouldReturn404(t *testing.T) {
	startContentAPIMock("notFound")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/content/5c3cae78-dbef-11e6-9d7c-be108f1c1dce")
	if err != nil {
		assert.FailNow(t, "Cannot send request to content endpoint", err.Error())
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Response status should be 404")
}


func TestShouldReturn400WhenInvalidUUID(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/content/123-invalid-uuid")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Response status should be 400")
}

func TestShouldReturn503WhenContentUnavailable(t *testing.T) {
	startContentAPIMock("unHappy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(contentAPIMock.URL + "/content/5c3cae78-dbef-11e6-9d7c-be108f1c1dce")
	if err != nil {
		assert.FailNow(t, "Cannot send request to content endpoint", err.Error())
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}


func TestShouldBeHealthy(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__health")
	if err != nil {
		assert.FailNow(t, "Cannot send request to health endpoint", err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeHealthyWhenContentApiIsNotHappy(t *testing.T) {
	startContentAPIMock("unhappy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(contentAPIMock.URL + "/__health")
	if err != nil {
		assert.FailNow(t, "Cannot send request to health endpoint", err.Error())
	}

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}


func TestShouldBeGoodToGo(t *testing.T) {
	startContentAPIMock("happy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(imageResolver.URL + "/__gtg")
	if err != nil {
		assert.FailNow(t, "Cannot send request to gtg endpoint", err.Error())
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeGoodToGoWhenContentApiIsNotHappy(t *testing.T) {
	startContentAPIMock("unhappy")
	startImageResolverService()
	defer stopServices()

	resp, err := http.Get(contentAPIMock.URL + "/__gtg")
	if err != nil {
		assert.FailNow(t, "Cannot send request to gtg endpoint", err.Error())
	}

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}