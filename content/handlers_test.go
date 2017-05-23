package content

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"bytes"
	"strings"
	"github.com/pkg/errors"
)

type ImageResolverMock struct {
	mockUnrollImages     func(c Content) (Content, error)
	mockUnrollLeadImages func(c Content) (Content, error)
}

func (irm *ImageResolverMock) UnrollImages(c Content) (Content, error) {
	return irm.mockUnrollImages(c)
}

func (irm *ImageResolverMock) UnrollLeadImages(c Content) (Content, error) {
	return irm.mockUnrollLeadImages(c)
}

func TestContentHandler_GetContentImages(t *testing.T) {
	ir := ImageResolverMock{
		mockUnrollImages: func(c Content) (Content, error) {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/valid-expanded-content-response.json")
			assert.NoError(t, err, "Cannot read resources test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return r, nil
		},
	}

	h := Handler{
		Service: &ir,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content/image", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentImages)

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

	req, err := http.NewRequest(http.MethodPost, "/content/image", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentImages)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestContentHandler_GetContentImages_InternalServerErrorWhenServiceReturnsError(t *testing.T) {
	ir := ImageResolverMock{
		mockUnrollImages: func(c Content) (Content, error) {
			return nil, errors.New("Image resolver failed")
		},
	}

	h := Handler{
		Service: &ir,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content/image", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentImages)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestContentHandler_GetLeadImages(t *testing.T) {
	ir := ImageResolverMock{
		mockUnrollLeadImages: func(c Content) (Content, error) {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/valid-expanded-internalcontent-response.json")
			assert.NoError(t, err, "Cannot read test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return r, nil
		},
	}

	h := Handler{
		Service: &ir,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent/image", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetLeadImages)

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

	req, err := http.NewRequest(http.MethodPost, "/internalcontent/image", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetLeadImages)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_GetLeadImages_InternalServerErrorWhenServiceReturnsError(t *testing.T) {
	ir := ImageResolverMock{
		mockUnrollLeadImages: func(c Content) (Content, error) {
			return nil, errors.New("Image resolver failed")
		},
	}

	h := Handler{
		Service: &ir,
	}

	body, err := ioutil.ReadFile("../test-resources/valid-article-internalcontent.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent/image", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetLeadImages)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

}
