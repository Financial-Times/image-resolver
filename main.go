package main

import (
	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Financial-Times/image-resolver/content"
	"github.com/Financial-Times/service-status-go/gtg"
	"github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"net"
	"net/http"
	"os"
	"time"
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
		Name:   "Content Source Application Name",
		Value:  "content-public-read",
		Desc:   "Content read app",
		EnvVar: "CONTENT_SOURCE_APP_NAME",
	})
	contentSourceURL := app.String(cli.StringOpt{
		Name:   "Content Source URL",
		Value:  "http://localhost:8080",
		Desc:   "Vulcan host",
		EnvVar: "CONTENT_SOURCE_URL",
	})
	contentSourcePath := app.String(cli.StringOpt{
		Name:   "Content Source Path",
		Value:  "/content",
		Desc:   "Endpoint of content source app for retrieving content",
		EnvVar: "CONTENT_SOURCE_PATH",
	})
	graphiteTCPAddress := app.String(cli.StringOpt{
		Name:   "graphite-tcp-address",
		Desc:   "Graphite TCP address, e.g. graphite.ft.com:2003. Leave as default if you do NOT want to output to graphite (e.g. if running locally)",
		EnvVar: "GRAPHITE_TCP_ADDRESS",
	})
	graphitePrefix := app.String(cli.StringOpt{
		Name:   "graphite-prefix",
		Value:  "coco.services.$ENV.image-resolver.0",
		Desc:   "Prefix to use. Should start with content, include the environment, and the host name. e.g. coco.pre-prod.image-resolver.1",
		EnvVar: "GRAPHITE_PREFIX",
	})
	logMetrics := app.Bool(cli.BoolOpt{
		Name:   "log-metrics",
		Value:  false,
		Desc:   "Whether to log metrics. Set to true if running locally and you want metrics output",
		EnvVar: "LOG_METRICS",
	})
	embeddedContentTypeWhitelist := app.String(cli.StringOpt{
		Name:   "embeddedContentTypeWhitelist",
		Value:  "http://www.ft.com/ontology/content/ImageSet",
		Desc:   "The type suported for embedded images, ex ImageSet",
		EnvVar: "EMBEDS_CONTENT_TYPE_WHITELIST",
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
			ContentSourceURL:     *contentSourceURL,
			HttpClient:           httpClient,
		}

		reader := content.NewContentReader(*contentSourceAppName, *contentSourceURL, *contentSourcePath, httpClient)
		parser := content.NewBodyParser(*embeddedContentTypeWhitelist)
		ir := content.NewImageResolver(reader, parser, *apiHost)

		baseftrwapp.OutputMetricsIfRequired(*graphiteTCPAddress, *graphitePrefix, *logMetrics)
		ch := &content.Handler{Service: ir}
		h := setupServiceHandler(sc, ch)
		err := http.ListenAndServe(":" + *port, h)
		if err != nil {
			log.Fatalf("Unable to start server: %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func setupServiceHandler(sc content.ServiceConfig, ch *content.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content/image", ch.GetContentImages).Methods("POST")
	r.HandleFunc("/internalcontent/image", ch.GetLeadImages).Methods("POST")

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	checks := []fthealth.Check{sc.ContentCheck()}
	hc := fthealth.HealthCheck{SystemCode: AppCode, Name: AppName, Description: AppDesc, Checks: checks}
	r.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(fthealth.Handler(&hc))})
	gtgHandler := httphandlers.NewGoodToGoHandler(gtg.StatusChecker(sc.GtgCheck))
	r.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(gtgHandler)})
	return r
}
