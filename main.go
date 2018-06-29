package main

import (
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
	AppDesc = "Image Resolver - unroll images for a given content"
)

func main() {
	app := cli.App(AppCode, AppDesc)
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "9090",
		Desc:   "application port",
		EnvVar: "PORT",
	})
	contentSourceAppName := app.String(cli.StringOpt{
		Name:   "contentSourceApplicationName",
		Value:  "content-public-read",
		Desc:   "Content read app",
		EnvVar: "CONTENT_SOURCE_APP_NAME",
	})
	contentSourceURL := app.String(cli.StringOpt{
		Name:   "contentSourceURL",
		Value:  "http://localhost:8080/__content-public-read/content",
		Desc:   "URL of the content source app",
		EnvVar: "CONTENT_SOURCE_URL",
	})
	contentSourceHealthURL := app.String(cli.StringOpt{
		Name:   "contentSourceHealthURL",
		Value:  "http://localhost:8080/__content-public-read/__health",
		Desc:   "Health url of the content source app",
		EnvVar: "CONTENT_SOURCE_HEALTH_URL",
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

	internalContentSourceAppName := app.String(cli.StringOpt{
		Name:   "internalComponentsSourceAppName",
		Value:  "document-store-api",
		Desc:   "Document Store API app",
		EnvVar: "INTERNAL_CONTENT_SOURCE_APP_NAME",
	})
	internalContentSourceURL := app.String(cli.StringOpt{
		Name:   "internalComponentsSourceURL",
		Value:  "http://localhost:8080/__document-store-api/internalcomponents",
		Desc:   "URL of Document Store API app",
		EnvVar: "INTERNAL_CONTENT_SOURCE_URL",
	})
	/*internalContentSourceHealthURL := app.String(cli.StringOpt{
		Name:   "internalComponentsSourceHealthURL",
		Value:  "http://document-store-api:8080/__health",
		Desc:   "Health URL of Document Store API app",
		EnvVar: "INTERNAL_CONTENT_SOURCE_HEALTH_URL",
	})*/

	nativeContentSourceAppName := app.String(cli.StringOpt{
		Name:   "nativeContentSourceAppName",
		Value:  "methode-api",
		Desc:   "Service name of the Native Content Source Application endpoint",
		EnvVar: "NATIVE_CONTENT_SOURCE_APP_NAME",
	})
	nativeContentSourceAppURL := app.String(cli.StringOpt{
		Name:   "nativeContentSourceAppURL",
		Value:  "http://methode-api-uk-t.svc.ft.com/eom-file/",
		Desc:   "URI of the Native Content Source Application endpoint",
		EnvVar: "NATIVE_CONTENT_SOURCE_APP_URL",
	})
	nativeContentSourceAppAuth := app.String(cli.StringOpt{
		Name:   "nativeContentSourceAppAuth",
		Value:  "default",
		Desc:   "Basic auth for Native Content Source Application",
		EnvVar: "NATIVE_CONTENT_SOURCE_APP_AUTH",
	})
	transformContentSourceURL := app.String(cli.StringOpt{
		Name:   "transformContentSourceURL",
		Value:  "http://localhost:8080/__methode-article-mapper/map",
		Desc:   "Methode Article Mapper URL",
		EnvVar: "TRANSFORM_CONTENT_SOURCE_APP_URL",
	})
	transformContentSourceAppName := app.String(cli.StringOpt{
		Name:   "transformContentSourceAppName",
		Value:  "methode-article-mapper",
		Desc:   "Methode Article Mapper app",
		EnvVar: "TRANSFORM_CONTENT_SOURCE_APP_NAME",
	})
	transformInternalContentSourceURL := app.String(cli.StringOpt{
		Name:   "transformContentSourceURL",
		Value:  "http://localhost:8080/__methode-article-internal-components-mapper/map",
		Desc:   "Methode Article Mapper URL",
		EnvVar: "TRANSFORM_INTERNAL_CONTENT_SOURCE_APP_URL",
	})
	transformInternalContentSourceAppName := app.String(cli.StringOpt{
		Name:   "transformContentSourceAppName",
		Value:  "methode-article-internal-components-mapper",
		Desc:   "Methode Article Mapper app",
		EnvVar: "TRANSFORM_INTERNAL_CONTENT_SOURCE_APP_NAME",
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
				Dial: (&net.Dialer{
					KeepAlive: 30 * time.Second,
				}).Dial,
			},
		}

		sc := content.ServiceConfig{
			ContentSourceAppName: *contentSourceAppName,
			ContentSourceURL:     *contentSourceHealthURL,
			HttpClient:           httpClient,
		}

		readerConfig := content.ReaderConfig{
			ContentSourceAppName:                  *contentSourceAppName,
			ContentSourceAppURL:                   *contentSourceURL,
			InternalContentSourceAppName:          *internalContentSourceAppName,
			InternalContentSourceAppURL:           *internalContentSourceURL,
			NativeContentSourceAppName:            *nativeContentSourceAppName,
			NativeContentSourceAppURL:             *nativeContentSourceAppURL,
			NativeContentSourceAppAuth:            *nativeContentSourceAppAuth,
			TransformContentSourceURL:             *transformContentSourceURL,
			TransformContentSourceAppName:         *transformContentSourceAppName,
			TransformInternalContentSourceURL:     *transformInternalContentSourceURL,
			TransformInternalContentSourceAppName: *transformInternalContentSourceAppName,
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

func setupServiceHandler(s *content.ContentUnroller, sc content.ServiceConfig) *mux.Router {
	r := mux.NewRouter()
	ch := &content.Handler{Service: s}
	r.HandleFunc("/content", ch.GetContent).Methods("POST")
	r.HandleFunc("/content-preview", ch.GetContentPreview).Methods("POST")
	r.HandleFunc("/internalcontent", ch.GetInternalContent).Methods("POST")
	r.HandleFunc("/internalcontent-preview", ch.GetInternalContentPreview).Methods("POST")

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	checks := []fthealth.Check{sc.ContentCheck()}
	hc := fthealth.TimedHealthCheck{
		HealthCheck: fthealth.HealthCheck{SystemCode: AppCode, Name: AppName, Description: AppDesc, Checks: checks},
		Timeout:     10 * time.Second,
	}

	r.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(fthealth.Handler(&hc))})
	gtgHandler := httphandlers.NewGoodToGoHandler(gtg.StatusChecker(sc.GtgCheck))
	r.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(gtgHandler)})
	return r
}
