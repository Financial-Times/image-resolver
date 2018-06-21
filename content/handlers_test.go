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

type ImageResolverMock struct {
	mockUnrollImages          func(req UnrollEvent) UnrollResult
	mockUnrollInternalContent func(req UnrollEvent) UnrollResult
}

func (irm *ImageResolverMock) UnrollImages(req UnrollEvent) UnrollResult {
	return irm.mockUnrollImages(req)
}

func (irm *ImageResolverMock) UnrollInternalContent(req UnrollEvent) UnrollResult {
	return irm.mockUnrollInternalContent(req)
}

func TestContentHandler_GetContentImages(t *testing.T) {
	ir := ImageResolverMock{
		mockUnrollImages: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/valid-expanded-content-response.json")
			assert.NoError(t, err, "Cannot read resources test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{
		Service: &ir,
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
	ir := ImageResolverMock{
		mockUnrollImages: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Image resolver failed")}
		},
	}

	h := Handler{
		Service: &ir,
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
	ir := ImageResolverMock{
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
		Service: &ir,
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
	ir := ImageResolverMock{
		mockUnrollInternalContent: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Image resolver failed")}
		},
	}

	h := Handler{
		Service: &ir,
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
