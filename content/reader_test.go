package content

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = []string{
	"639cd952-149f-11e7-2ea7-a07ecd9ac73f",
	"639cd952-149f-11e7-2ea7-a07ecd9ac73f",
	"71231d3a-13c7-11e7-2ea7-a07ecd9ac73f",
	"d02886fc-58ff-11e8-9859-6668838a4c10",
}

var dynamicContentTestData = []string{"d02886fc-58ff-11e8-9859-6668838a4c10"}

func successfulContentServerMock(t *testing.T, resource string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(resource)
		if err != nil {
			assert.NoError(t, err, "File necessary for starting mock server not found.")
			return
		}
		defer file.Close()
		io.Copy(w, file)
	}))
}

func errorContentServerMock(t *testing.T, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotEqual(t, http.StatusOK, statusCode, fmt.Sprintf("Status code should not be %d", http.StatusOK))
		w.WriteHeader(statusCode)
	}))
}

func readerForTest(cfg ReaderConfig) *ContentReader {
	return &ContentReader{
		config: cfg,
		client: http.DefaultClient,
	}
}

func TestGet(t *testing.T) {
	ts := successfulContentServerMock(t, "../test-resources/source-content-valid-response.json")
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	body, err := ioutil.ReadFile("../test-resources/reader-content-valid-response.json")
	assert.NoError(t, err, "Cannot read file necessary for running test case.")
	var expected map[string]Content
	err = json.Unmarshal(body, &expected)
	assert.NoError(t, err, "Cannot read expected response for test case.")

	actual, err := cr.Get(testData, "tid_1")
	assert.NoError(t, err, "Error while getting content data")
	assert.Equal(t, expected, actual)
}

func TestGet_ContentSourceReturns500(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusInternalServerError)
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.Get(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGet_ContentSourceReturns404(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusNotFound)
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.Get(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGet_ContentSourceCannotBeResolved(t *testing.T) {
	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    "http://sampleAddress:8080/content",
	}
	cr := readerForTest(cfg)

	_, err := cr.Get(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGet_ContentSourceHasInvalidURL(t *testing.T) {
	cfg := ReaderConfig{
		ContentStoreAppName: "&&^%&&^",
		ContentStoreHost:    "@$@",
	}
	cr := readerForTest(cfg)

	_, err := cr.Get(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternal(t *testing.T) {
	ts := successfulContentServerMock(t, "../test-resources/internalcontent-source-valid-response.json")
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	body, err := ioutil.ReadFile("../test-resources/reader-internalcontent-dynamic-valid-response.json")
	assert.NoError(t, err, "Cannot read file necessary for running test case.")
	var expected map[string]Content
	err = json.Unmarshal(body, &expected)
	assert.NoError(t, err, "Cannot read expected response for test case.")

	actual, err := cr.GetInternal(testData, "tid_1")
	assert.NoError(t, err, "Error while getting content data")
	assert.Equal(t, expected, actual)
}

func TestGetInternal_ContentSourceReturns500(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusInternalServerError)
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternal(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternal_ContentSourceReturns404(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusNotFound)
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternal(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternal_ContentSourceCannotBeResolved(t *testing.T) {
	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    "http://sampleAddress:8080/content",
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternal(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternal_ContentSourceHasInvalidURL(t *testing.T) {
	cfg := ReaderConfig{
		ContentStoreAppName: "&&^%&&^",
		ContentStoreHost:    "@$@",
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternal(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetPreview(t *testing.T) {
	ts := successfulContentServerMock(t, "../test-resources/source-contentpreview-valid-response.json")
	defer ts.Close()

	cfg := ReaderConfig{
		ContentPreviewAppName: "content-preview-source-app-name",
		ContentPreviewHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	body, err := ioutil.ReadFile("../test-resources/reader-contentpreview-dynamic-content-valid-response.json")
	assert.NoError(t, err, "Cannot read file necessary for running test case.")
	var expected map[string]Content
	err = json.Unmarshal(body, &expected)
	assert.NoError(t, err, "Cannot read expected response for test case.")

	actual, err := cr.GetPreview(dynamicContentTestData, "tid_1")
	assert.NoError(t, err, "Error while getting content data")
	assert.Equal(t, expected, actual)
}

func TestGetPreview_ContentSourceReturns500(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusInternalServerError)
	defer ts.Close()

	cfg := ReaderConfig{
		ContentPreviewAppName: "content-preview-source-app-name",
		ContentPreviewHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.GetPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetPreview_ContentSourceReturns404(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusNotFound)
	defer ts.Close()

	cfg := ReaderConfig{
		ContentPreviewAppName: "content-preview-source-app-name",
		ContentPreviewHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.GetPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetPreview_ContentSourceCannotBeResolved(t *testing.T) {
	cfg := ReaderConfig{
		ContentPreviewAppName: "content-preview-source-app-name",
		ContentPreviewHost:    "http://sampleAddress:8080/content",
	}
	cr := readerForTest(cfg)

	_, err := cr.GetPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetPreview_ContentSourceHasInvalidURL(t *testing.T) {
	cfg := ReaderConfig{
		ContentPreviewAppName: "&&^%&&^",
		ContentPreviewHost:    "@$@",
	}
	cr := readerForTest(cfg)

	_, err := cr.GetPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternalPreview(t *testing.T) {
	ts := successfulContentServerMock(t, "../test-resources/source-internalcontentpreview-valid-response.json")
	defer ts.Close()

	cfg := ReaderConfig{
		InternalContentPreviewAppName: "internal-content-preview-source-app-name",
		InternalContentPreviewHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	body, err := ioutil.ReadFile("../test-resources/reader-internalcontentpreview-dynamic-content-valid-response.json")
	assert.NoError(t, err, "Cannot read file necessary for running test case.")
	var expected map[string]Content
	err = json.Unmarshal(body, &expected)
	assert.NoError(t, err, "Cannot read expected response for test case.")

	actual, err := cr.GetInternalPreview(dynamicContentTestData, "tid_1")
	assert.NoError(t, err, "Error while getting content data")
	assert.Equal(t, expected, actual)
}

func TestGetInternalPreview_ContentSourceReturns500(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusInternalServerError)
	defer ts.Close()

	cfg := ReaderConfig{
		InternalContentPreviewAppName: "internal-content-preview-source-app-name",
		InternalContentPreviewHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternalPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternalPreview_ContentSourceReturns404(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusNotFound)
	defer ts.Close()

	cfg := ReaderConfig{
		InternalContentPreviewAppName: "internal-content-preview-source-app-name",
		InternalContentPreviewHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternalPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternalPreview_ContentSourceCannotBeResolved(t *testing.T) {
	cfg := ReaderConfig{
		InternalContentPreviewAppName: "internal-content-preview-source-app-name",
		InternalContentPreviewHost:    "http://sampleAddress:8080/content",
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternalPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestGetInternalPreview_ContentSourceHasInvalidURL(t *testing.T) {
	cfg := ReaderConfig{
		InternalContentPreviewAppName: "&&^%&&^",
		InternalContentPreviewHost:    "@$@",
	}
	cr := readerForTest(cfg)

	_, err := cr.GetInternalPreview(dynamicContentTestData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}
