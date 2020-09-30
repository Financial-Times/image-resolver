package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Financial-Times/content-unroller/content"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/service-status-go/gtg"
	"github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

const (
	AppCode = "content-unroller"
	AppName = "Content Unroller"
	AppDesc = "Content Unroller - unroll images and dynamic content for a given content"
)

func main() {
	app := cli.App(AppCode, AppDesc)
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "9090",
		Desc:   "application port",
		EnvVar: "PORT",
	})
	contentStoreApplicationName := app.String(cli.StringOpt{
		Name:   "contentSourceAppName",
		Value:  "content-public-read",
		Desc:   "Content read app",
		EnvVar: "CONTENT_STORE_APP_NAME",
	})
	contentStoreHost := app.String(cli.StringOpt{
		Name:   "contentStoreHost",
		Value:  "http://localhost:8080/__content-public-read",
		Desc:   "Content source hostname",
		EnvVar: "CONTENT_STORE_HOST",
	})
	contentPreviewAppName := app.String(cli.StringOpt{
		Name:   "contentPreviewAppName",
		Value:  "content-public-read-preview",
		Desc:   "Content Preview app",
		EnvVar: "CONTENT_PREVIEW_APP_NAME",
	})
	contentPreviewHost := app.String(cli.StringOpt{
		Name:   "contentPreviewHost",
		Value:  "http://localhost:8080/__content-preview",
		Desc:   "Content Preview hostname",
		EnvVar: "CONTENT_PREVIEW_HOST",
	})
	contentPathEndpoint := app.String(cli.StringOpt{
		Name:   "contentPathEndpoint",
		Value:  "/content",
		Desc:   "/content path",
		EnvVar: "CONTENT_PATH",
	})
	internalContentPathEndpoint := app.String(cli.StringOpt{
		Name:   "internalContentPathEndpoint",
		Value:  "/internalcontent",
		Desc:   "/internalcontent path",
		EnvVar: "INTERNAL_CONTENT_PATH",
	})
	apiHost := app.String(cli.StringOpt{
		Name:   "apiHost",
		Value:  "test.api.ft.com",
		Desc:   "API host to use for URLs in responses",
		EnvVar: "API_HOST",
	})

	app.Action = func() {
		httpClient := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				DialContext: (&net.Dialer{
					KeepAlive: 30 * time.Second,
				}).DialContext,
			},
		}

		sc := content.ServiceConfig{
			ContentStoreAppName:        *contentStoreApplicationName,
			ContentStoreAppHealthURI:   getServiceHealthURI(*contentStoreHost),
			ContentPreviewAppName:      *contentPreviewAppName,
			ContentPreviewAppHealthURI: getServiceHealthURI(*contentPreviewHost),
			HTTPClient:                 httpClient,
		}

		readerConfig := content.ReaderConfig{
			ContentStoreAppName:         *contentStoreApplicationName,
			ContentStoreHost:            *contentStoreHost,
			ContentPreviewAppName:       *contentPreviewAppName,
			ContentPreviewHost:          *contentPreviewHost,
			ContentPathEndpoint:         *contentPathEndpoint,
			InternalContentPathEndpoint: *internalContentPathEndpoint,
		}

		reader := content.NewContentReader(readerConfig, httpClient)
		unroller := content.NewContentUnroller(reader, *apiHost)

		h := setupServiceHandler(unroller, sc)
		err := http.ListenAndServe(":"+*port, h)
		if err != nil {
			log.Fatalf("Unable to start server: %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func setupServiceHandler(s content.Unroller, sc content.ServiceConfig) *mux.Router {
	r := mux.NewRouter()
	ch := &content.Handler{Service: s}

	var checks []fthealth.Check
	var gtgHandler func(http.ResponseWriter, *http.Request)

	r.HandleFunc("/content", ch.GetContent).Methods("POST")
	r.HandleFunc("/internalcontent", ch.GetInternalContent).Methods("POST")
	checks = []fthealth.Check{sc.ContentStoreCheck()}
	gtgHandler = httphandlers.NewGoodToGoHandler(gtg.StatusChecker(sc.GtgCheck))

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	hc := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{SystemCode: AppCode, Name: AppName, Description: AppDesc, Checks: checks},
		Timeout:     10 * time.Second,
	}

	r.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(fthealth.Handler(&hc))})
	r.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(gtgHandler)})
	return r
}

func getServiceHealthURI(hostname string) string {
	return fmt.Sprintf("%s%s", hostname, "/__health")
}
