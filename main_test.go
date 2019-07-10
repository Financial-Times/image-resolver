package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/Financial-Times/content-unroller/content"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	contentStoreAppName           = "content-source-app-name"
	contentPreviewAppName         = "content-preview-app-name"
	internalContentPreviewAppName = "internal-content-preview-app-name"
)

var (
	unrollerService *httptest.Server
)

func startContentServerMock(resource string, isPreview bool) *httptest.Server {
	router := mux.NewRouter()
	router.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(statusOkHandler)})
	router.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(statusOkHandler)})
	router.Path("/").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(successfulContentServerMock(resource))})
	if isPreview {
		router.Path("/{uuid}").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(successfulContentServerMock(resource))})
	}

	return httptest.NewServer(router)
}

func successfulContentServerMock(resource string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(resource)
		if err != nil {
			return
		}
		defer file.Close()
		io.Copy(w, file)
	}
}

func statusOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func startUnhealthyContentServerMock() *httptest.Server {
	router := mux.NewRouter()
	router.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(statusServiceUnavailableHandler)})
	router.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(statusServiceUnavailableHandler)})

	return httptest.NewServer(router)
}

func statusServiceUnavailableHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func startUnrollerService(contentStoreURL string, contentPreviewURL string, flow string) {
	sc := content.ServiceConfig{
		ContentStoreAppName:        contentStoreAppName,
		ContentStoreAppHealthURI:   getServiceHealthURI(contentStoreURL),
		ContentPreviewAppName:      contentPreviewAppName,
		ContentPreviewAppHealthURI: getServiceHealthURI(contentPreviewURL),
		HTTPClient:                 http.DefaultClient,
	}

	rc := content.ReaderConfig{
		ContentStoreAppName:         contentStoreAppName,
		ContentStoreHost:            contentStoreURL,
		ContentPreviewAppName:       contentPreviewAppName,
		ContentPreviewHost:          contentPreviewURL,
		ContentPathEndpoint:         "",
		InternalContentPathEndpoint: "",
	}

	reader := content.NewContentReader(rc, http.DefaultClient)
	unroller := content.NewContentUnroller(reader, "test.api.ft.com")

	h := setupServiceHandler(unroller, sc, flow)
	unrollerService = httptest.NewServer(h)
}
