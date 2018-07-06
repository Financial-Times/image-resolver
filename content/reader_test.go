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

const testResourcesRoot = "../test-resources/"

var testData = []string{
	"639cd952-149f-11e7-2ea7-a07ecd9ac73f",
	"639cd952-149f-11e7-2ea7-a07ecd9ac73f",
	"71231d3a-13c7-11e7-2ea7-a07ecd9ac73f",
}

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

func TestContentReader_Get(t *testing.T) {
	ts := successfulContentServerMock(t, testResourcesRoot+"valid-content-source-response.json")
	defer ts.Close()

	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    ts.URL,
	}
	cr := readerForTest(cfg)

	b, err := ioutil.ReadFile(testResourcesRoot + "valid-content-reader-response.json")
	assert.NoError(t, err, "Cannot read file necessary for running test case.")
	var expected map[string]Content
	err = json.Unmarshal(b, &expected)
	assert.NoError(t, err, "Cannot read expected response for test case.")

	actual, err := cr.Get(testData, "tid_1")
	assert.NoError(t, err, "Error while getting content data")
	assert.Equal(t, expected, actual)
}

func TestContentReader_Get_ContentSourceReturns500(t *testing.T) {
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

func TestContentReader_Get_ContentSourceReturns404(t *testing.T) {
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

func TestContentReader_Get_ContentSourceCannotBeResolved(t *testing.T) {
	cfg := ReaderConfig{
		ContentStoreAppName: "content-source-app-name",
		ContentStoreHost:    "http://sampleAddress:8080/content",
	}
	cr := readerForTest(cfg)

	_, err := cr.Get(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}

func TestContentReader_Get_ContentSourceHasInvalidURL(t *testing.T) {
	cfg := ReaderConfig{
		ContentStoreAppName: "&&^%&&^",
		ContentStoreHost:    "@$@",
	}
	cr := readerForTest(cfg)

	_, err := cr.Get(testData, "tid_1")
	assert.Error(t, err, "There should an error thrown")
}
