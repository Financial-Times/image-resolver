package content

import (
	"testing"
	"net/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"encoding/json"
	"net/http/httptest"
	"os"
	"io"
	"fmt"
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

func readerForTest(URL string, path string) *ContentReader {
	return &ContentReader{
		contentAppName: "content-source-app-name",
		contentAppURL:  URL,
		path:           path,
		client:         http.DefaultClient,
	}
}

func TestContentReader_Get(t *testing.T) {
	ts := successfulContentServerMock(t, testResourcesRoot + "valid-content-source-response.json")
	defer ts.Close()

	cr := readerForTest(ts.URL, "/content")

	b, err := ioutil.ReadFile(testResourcesRoot + "valid-content-reader-response.json")
	assert.NoError(t, err, "Cannot read file necessary for running test case.")
	var expected map[string]Content
	err = json.Unmarshal(b, &expected)
	assert.NoError(t, err, "Cannot read expected response for test case.")

	actual, err := cr.Get(testData)
	assert.NoError(t, err, "Error while getting content data")
	assert.Equal(t, expected, actual)
}

func TestContentReader_Get_ContentSourceReturns500(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusInternalServerError)
	defer ts.Close()

	cr := readerForTest(ts.URL, "/content")

	_, err := cr.Get(testData)
	assert.Error(t, err, "There should an error thrown")
}

func TestContentReader_Get_ContentSourceReturns404(t *testing.T) {
	ts := errorContentServerMock(t, http.StatusNotFound)
	defer ts.Close()

	cr := readerForTest(ts.URL, "/content")

	_, err := cr.Get(testData)
	assert.Error(t, err, "There should an error thrown")
}

func TestContentReader_Get_ContentSourceCannotBeResolved(t *testing.T) {
	cr := readerForTest("http://sampleAddress:8080", "/content")

	_, err := cr.Get(testData)
	assert.Error(t, err, "There should an error thrown")
}

func TestContentReader_Get_ContentSourceHasInvalidURL(t *testing.T) {
	cr := readerForTest("&&^%&&^", "@$@")

	_, err := cr.Get(testData)
	assert.Error(t, err, "There should an error thrown")
}

func TestNewContentReader(t *testing.T) {
	actual := NewContentReader("content-source-app", "http://localhost:8080", "/content", http.DefaultClient)
	expected := &ContentReader{
		contentAppName: "content-source-app",
		contentAppURL:  "http://localhost:8080",
		path:           "/content",
		client:         http.DefaultClient,
	}
	assert.Equal(t, expected, actual)
}
