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

const (
	InvalidBodyRequest = "{\"id\": \"d02886fc-58ff-11e8-9859-6668838a4c10\"}"
)

type ContentUnrollerMock struct {
	mockUnrollContent                func(UnrollEvent) UnrollResult
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

func TestGetContentReturns200(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollContent: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/content-valid-response.json")
			assert.NoError(t, err, "Cannot read resources test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/content-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, err := ioutil.ReadFile("../test-resources/content-valid-response.json")
	assert.NoError(t, err, "Cannot read test file")
	actualBody := rr.Body

	assert.JSONEq(t, string(expectedBody), string(actualBody.Bytes()))
}

func TestGetContent_UnrollEventError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/content", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "invalid character")
}

func TestGetContent_ValidationError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/content", strings.NewReader(InvalidBodyRequest))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Invalid content")
}

func TestGetContent_UnrollingError(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollContent: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Error while unrolling content")}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/content-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Error while unrolling content")
}

func TestGetInternalContentReturns200(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollInternalContent: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/internalcontent-valid-response.json")
			assert.NoError(t, err, "Cannot read test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/internalcontent-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, err := ioutil.ReadFile("../test-resources/internalcontent-valid-response.json")
	assert.NoError(t, err, "Cannot read test file")
	actualBody := rr.Body

	assert.JSONEq(t, string(expectedBody), string(actualBody.Bytes()))
}

func TestGetInternalContent_UnrollEventError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/internalcontent", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "invalid character")
}

func TestGetInternalContent_ValidationError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/internalcontent", strings.NewReader(InvalidBodyRequest))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Invalid content")
}

func TestGetInternalContent_UnrollingError(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollInternalContent: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Error while unrolling content")}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/internalcontent-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Error while unrolling content")
}
func TestGetContentPreviewReturns200(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollContentPreview: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/contentpreview-valid-response.json")
			assert.NoError(t, err, "Cannot read resources test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/content-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content-preview", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, err := ioutil.ReadFile("../test-resources/contentpreview-valid-response.json")
	assert.NoError(t, err, "Cannot read test file")
	actualBody := rr.Body

	assert.JSONEq(t, string(expectedBody), string(actualBody.Bytes()))
}
func TestGetContentPreview_UnrollEventError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/content-preview", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "invalid character")
}

func TestGetContentPreview_ValidationError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/content-preview", strings.NewReader(InvalidBodyRequest))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Invalid content")
}

func TestGetContentPreview_UnrollingError(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollContentPreview: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Error while unrolling content")}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/content-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/content-preview", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Error while unrolling content")
}

func TestGetInternalContentPreviewReturns200(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollInternalContentPreview: func(req UnrollEvent) UnrollResult {
			var r Content
			fileBytes, err := ioutil.ReadFile("../test-resources/internalcontentpreview-valid-response.json")
			assert.NoError(t, err, "Cannot read test file")
			err = json.Unmarshal(fileBytes, &r)
			assert.NoError(t, err, "Cannot build json body")
			return UnrollResult{r, nil}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/internalcontent-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent-preview", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	expectedBody, err := ioutil.ReadFile("../test-resources/internalcontentpreview-valid-response.json")
	assert.NoError(t, err, "Cannot read test file")
	actualBody := rr.Body

	assert.JSONEq(t, string(expectedBody), string(actualBody.Bytes()))
}
func TestGetInternalContentPreview_UnrollEventError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/internalcontent-preview", strings.NewReader("sample body"))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContent)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "invalid character")
}

func TestGetInternalContentPreview_ValidationError(t *testing.T) {
	h := Handler{nil}
	req, err := http.NewRequest(http.MethodPost, "/internalcontent-preview", strings.NewReader(InvalidBodyRequest))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Invalid content")
}

func TestGetInternalContentPreview_UnrollingError(t *testing.T) {
	cu := ContentUnrollerMock{
		mockUnrollInternalContentPreview: func(req UnrollEvent) UnrollResult {
			return UnrollResult{nil, errors.New("Error while unrolling content")}
		},
	}

	h := Handler{&cu}
	body, err := ioutil.ReadFile("../test-resources/internalcontent-valid-request.json")
	assert.NoError(t, err, "Cannot read test file")
	req, err := http.NewRequest(http.MethodPost, "/internalcontent-preview", bytes.NewReader(body))
	assert.NoError(t, err, "Cannot create request necessary for test")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetInternalContentPreview)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, string(rr.Body.Bytes()), "Error while unrolling content")
}
