package content

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type ContentUnrollerMock struct {
	mockUnrollContent                func(req UnrollEvent) UnrollResult
	mockUnrollContentPreview         func(UnrollEvent) UnrollResult
	mockUnrollInternalContent        func(UnrollEvent) UnrollResult
	mockUnrollInternalContentPreview func(UnrollEvent) UnrollResult
}

func (cu *ContentUnrollerMock) UnrollContent(req UnrollEvent) UnrollResult {
	return cu.mockUnrollContent(req)
}

func (cu *ContentUnrollerMock) UnrollContentPreview(req UnrollEvent) UnrollResult {
	return cu.mockUnrollContentPreview(req)
}

func (cu *ContentUnrollerMock) UnrollInternalContent(req UnrollEvent) UnrollResult {
	return cu.mockUnrollInternalContent(req)
}

func (cu *ContentUnrollerMock) UnrollInternalContentPreview(req UnrollEvent) UnrollResult {
	return cu.mockUnrollInternalContentPreview(req)
}

func TestContentHandler_GetContentImages(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollContent: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/valid-expanded-content-response.json")
			assert.NoError(t, err, "Cannot read resources test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{
		Service: &cu,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, err := ioutil.ReadFile("../test-resources/valid-expanded-content-response.json")
	assert.NoError(t, err, "Cannot read test file")
	actualBody := rr.Body

	assert.JSONEq(t, string(expectedBody), string(actualBody.Bytes()))
}

func TestContentHandler_GetContentImagesWhenBodyNotValid(t *testing.T) {
	h := Handler{
		Service: nil,
	}

	req, err := http.NewRequest(http.MethodPost, "/content", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestContentHandler_GetContentImages_InternalServerErrorWhenServiceReturnsError(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollContent: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Image resolver failed")}
		},
	}

	h := Handler{
		Service: &cu,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestContentHandler_GetLeadImages(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollInternalContent: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/valid-expanded-internalcontent-response.json")
			assert.NoError(t, err, "Cannot read test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{
		Service: &cu,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, err := ioutil.ReadFile("../test-resources/valid-expanded-internalcontent-response.json")
	assert.NoError(t, err, "Cannot read test file")
	actualBody := rr.Body

	assert.JSONEq(t, string(expectedBody), string(actualBody.Bytes()))
}

func TestContentHandler_GetLeadImagesWhenBodyNotValid(t *testing.T) {
	h := Handler{
		Service: nil,
	}

	req, err := http.NewRequest(http.MethodPost, "/internalcontent", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_GetLeadImages_InternalServerErrorWhenServiceReturnsError(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollInternalContent: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Image resolver failed")}
		},
	}

	h := Handler{
		Service: &cu,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
