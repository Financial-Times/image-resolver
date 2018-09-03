package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/image-resolver/content"
	"github.com/Financial-Times/service-status-go/gtg"
	"github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

const (
	AppCode = "image-resolver"
	AppName = "Image Resolver"
	AppDesc = "Image Resolver - unroll images and dynamic content for a given content"
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
	graphiteTCPAddress := app.String(cli.StringOpt{
		Name:   "graphiteTCPAddress",
		Desc:   "Graphite TCP address, e.g. graphite.ft.com:2003. Leave as default if you do NOT want to output to graphite (e.g. if running locally)",
		EnvVar: "GRAPHITE_TCP_ADDRESS",
	})
	graphitePrefix := app.String(cli.StringOpt{
		Name:   "graphitePrefix",
		Value:  "coco.services.$ENV.image-resolver.0",
		Desc:   "Prefix to use. Should start with content, include the environment, and the host name. e.g. coco.pre-prod.image-resolver.1",
		EnvVar: "GRAPHITE_PREFIX",
	})
	logMetrics := app.Bool(cli.BoolOpt{
		Name:   "logMetrics",
		Value:  false,
		Desc:   "Whether to log metrics. Set to true if running locally and you want metrics output",
		EnvVar: "LOG_METRICS",
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

		baseftrwapp.OutputMetricsIfRequired(*graphiteTCPAddress, *graphitePrefix, *logMetrics)
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
	r.HandleFunc("/content", ch.GetContent).Methods("POST")
	r.HandleFunc("/content-preview", ch.GetContentPreview).Methods("POST")
	r.HandleFunc("/internalcontent", ch.GetInternalContent).Methods("POST")
	r.HandleFunc("/internalcontent-preview", ch.GetInternalContentPreview).Methods("POST")

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	checks := []fthealth.Check{sc.ContentStoreCheck(), sc.ContentPreviewCheck()}
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
