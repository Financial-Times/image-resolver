// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	)

func TestContent_ShouldReturn200(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	startUnrollerService(contentStoreServiceMock.URL, "")
	defer contentStoreServiceMock.Close()
	defer unrollerService.Close()

	expected, err := ioutil.ReadFile("test-resources/content-valid-response.json")
	assert.NoError(t, err, "")

	body, err := ioutil.ReadFile("test-resources/content-valid-request.json")
	assert.NoError(t, err, "Cannot read file necessary for test case")
	resp, err := http.Post(unrollerService.URL+"/content", "application/json", bytes.NewReader(body))
	assert.NoError(t, err, "Should not fail")
	defer resp.Body.Close()
	actualResponse, err := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, err, "")
	assert.JSONEq(t, string(expected), string(actualResponse))
}

func TestContent_ShouldReturn400WhenInvalidJson(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	startUnrollerService(contentStoreServiceMock.URL, "")
	defer contentStoreServiceMock.Close()
	defer unrollerService.Close()

	body := `{"body":"invalid""body"}`
	resp, err := http.Post(unrollerService.URL+"/content", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestContent_ShouldReturn400WhenInvalidContentRequest(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	startUnrollerService(contentStoreServiceMock.URL, "")
	defer contentStoreServiceMock.Close()
	defer unrollerService.Close()

	body := `{"id":"36037ab1-da3b-35bf-b5ee-4fc23723b635"}`
	resp, err := http.Post(unrollerService.URL+"/content", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestInternalContent_ShouldReturn200(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/internalcontent-source-valid-response.json", false)
	startUnrollerService(contentStoreServiceMock.URL, "")
	defer contentStoreServiceMock.Close()
	defer unrollerService.Close()

	expected, err := ioutil.ReadFile("test-resources/internalcontent-valid-response-no-lead-images.json")
	assert.NoError(t, err, "")

	body, err := ioutil.ReadFile("test-resources/internalcontent-valid-request.json")
	assert.NoError(t, err, "Cannot read file necessary for test case")
	resp, err := http.Post(unrollerService.URL+"/internalcontent", "application/json", bytes.NewReader(body))
	assert.NoError(t, err, "Should not fail")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	actualResponse, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "")
	assert.JSONEq(t, string(expected), string(actualResponse))
}

func TestInternalContent_ShouldReturn400InvalidJson(t *testing.T) {
	internalContentStoreServiceMock := startContentServerMock("test-resources/internalcontent-source-valid-response.json", false)
	startUnrollerService(internalContentStoreServiceMock.URL, "")
	defer internalContentStoreServiceMock.Close()
	defer unrollerService.Close()

	body := `{"body":"invalid""body"}`
	resp, err := http.Post(unrollerService.URL+"/internalcontent", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestInternalContent_ShouldReturn400InvalidArticle(t *testing.T) {
	internalContentStoreServiceMock := startContentServerMock("test-resources/internalcontent-source-valid-response.json", false)
	startUnrollerService(internalContentStoreServiceMock.URL, "")
	defer internalContentStoreServiceMock.Close()
	defer unrollerService.Close()

	body := `{"id":"36037ab1-da3b-35bf-b5ee-4fc23723b635"}`
	resp, err := http.Post(unrollerService.URL+"/internalcontent", "application/json", strings.NewReader(body))
	assert.NoError(t, err, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestShouldBeHealthy(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	contentPreviewServiceMock := startContentServerMock("test-resources/source-contentpreview-valid-response.json", true)
	startUnrollerService(contentStoreServiceMock.URL, contentPreviewServiceMock.URL)

	defer contentStoreServiceMock.Close()
	defer contentPreviewServiceMock.Close()
	defer unrollerService.Close()

	resp, err := http.Get(unrollerService.URL + "/__health")

	assert.NoError(t, err, "Cannot send request to health endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeHealthyWhenContentPreviewIsNotHappy(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	contentPreviewServiceMock := startUnhealthyContentServerMock()
	startUnrollerService(contentStoreServiceMock.URL, contentPreviewServiceMock.URL)

	defer contentStoreServiceMock.Close()
	defer contentPreviewServiceMock.Close()
	defer unrollerService.Close()

	resp, err := http.Get(unrollerService.URL + "/__health")

	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 503")

	respBody, err := ioutil.ReadAll(resp.Body)

	assert.NoError(t, err, "")
	assertMsg := fmt.Sprintf(`"id":"check-connect-%v","name":"Check connectivity to %v","ok":false`, contentPreviewAppName, contentPreviewAppName)
	assert.Contains(t, string(respBody), assertMsg)
}

func TestShouldBeGoodToGo(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	contentPreviewServiceMock := startContentServerMock("test-resources/source-contentpreview-valid-response.json", true)
	startUnrollerService(contentStoreServiceMock.URL, contentPreviewServiceMock.URL)

	defer contentStoreServiceMock.Close()
	defer contentPreviewServiceMock.Close()
	defer unrollerService.Close()

	resp, err := http.Get(unrollerService.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status should be 200")
}

func TestShouldNotBeGoodToGoWhenContentStoreIsNotHappy(t *testing.T) {
	contentStoreServiceMock := startUnhealthyContentServerMock()
	contentPreviewServiceMock := startContentServerMock("test-resources/source-contentpreview-valid-response.json", true)
	startUnrollerService(contentStoreServiceMock.URL, contentPreviewServiceMock.URL)

	defer contentStoreServiceMock.Close()
	defer contentPreviewServiceMock.Close()
	defer unrollerService.Close()

	resp, err := http.Get(unrollerService.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}

func TestShouldNotBeGoodToGoWhenContentPreviewIsNotHappy(t *testing.T) {
	contentStoreServiceMock := startContentServerMock("test-resources/source-content-valid-response.json", false)
	contentPreviewServiceMock := startUnhealthyContentServerMock()
	startUnrollerService(contentStoreServiceMock.URL, contentPreviewServiceMock.URL)

	defer contentStoreServiceMock.Close()
	defer contentPreviewServiceMock.Close()
	defer unrollerService.Close()

	resp, err := http.Get(unrollerService.URL + "/__gtg")
	assert.NoError(t, err, "Cannot send request to gtg endpoint")
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Response status should be 503")
}
