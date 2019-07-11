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
	"github.com/jawher/mow.cli"
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
	flow := app.String(cli.StringOpt{
		Name:   "flow",
		Value:  "read",
		Desc:   "Marks if it is read or preview flow (default: read)",
		EnvVar: "FLOW",
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

		switch *flow {
		case "read", "preview":
			log.Infof("Value of 'flow' set to '%s'.", *flow)
		default:
			log.Warnf("Value of 'flow' should be one of: 'read' or 'preview', defaulting to 'read'.")
		}

		h := setupServiceHandler(unroller, sc, *flow)
		err := http.ListenAndServe(":"+*port, h)
		if err != nil {
			log.Fatalf("Unable to start server: %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func setupServiceHandler(s content.Unroller, sc content.ServiceConfig, flow string) *mux.Router {
	r := mux.NewRouter()
	ch := &content.Handler{Service: s}
	// Splitting the read and preview flow: endpoints and healthchecks assigned accordingly
	var checks []fthealth.Check

	if flow == "preview" {
		r.HandleFunc("/content-preview", ch.GetContentPreview).Methods("POST")
		r.HandleFunc("/internalcontent-preview", ch.GetInternalContentPreview).Methods("POST")
		checks = []fthealth.Check{sc.ContentStoreCheck(), sc.ContentPreviewCheck()}
	} else {
		// the default value for flow is "read", there are no other cases other than "read" and "preview"
		// the purpose for such is setup is that the service can be run wothout specifiyng flow for easier testing (when a dependacy)
		r.HandleFunc("/content", ch.GetContent).Methods("POST")
		r.HandleFunc("/internalcontent", ch.GetInternalContent).Methods("POST")
		checks = []fthealth.Check{sc.ContentStoreCheck()}
	}

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	hc := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{SystemCode: AppCode, Name: AppName, Description: AppDesc, Checks: checks},
		Timeout:     10 * time.Second,
	}

	r.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(fthealth.Handler(&hc))})
	gtgHandler := httphandlers.NewGoodToGoHandler(gtg.StatusChecker(sc.GtgCheck))
	r.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(gtgHandler)})
	return r
}

func getServiceHealthURI(hostname string) string {
	return fmt.Sprintf("%s%s", hostname, "/__health")
}
