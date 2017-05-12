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
	AppCode        = "image-resolver"
	AppName        = "Image Resolver"
	AppDesc        = "Image Resolver - unroll images for a given content"
	Port           = "port"
	CprHost        = "cprHost"
	RouterAddress  = "routerAddress"
	EmbedsType     = "embedsType"
	GraphiteTCP    = "graphite-tcp-address"
	GraphitePrefix = "graphite-prefix"
	LogMetrics     = "log-metrics"
)

func main() {
	app := cli.App(AppCode, AppDesc)
	port := app.String(cli.StringOpt{
		Name:   Port,
		Value:  "9090",
		Desc:   "application port",
		EnvVar: "PORT"})
	cprHost := app.String(cli.StringOpt{
		Name:   CprHost,
		Value:  "content-public-read",
		Desc:   "Content read app",
		EnvVar: "CONTENT_READ_ENDPOINT",
	})
	routerAddress := app.String(cli.StringOpt{
		Name:   RouterAddress,
		Value:  "localhost:8080",
		Desc:   "Vulcan host",
		EnvVar: "ROUTER_ADDRESS",
	})
	graphiteTCPAddress := app.String(cli.StringOpt{
		Name:   GraphiteTCP,
		Desc:   "Graphite TCP address, e.g. graphite.ft.com:2003. Leave as default if you do NOT want to output to graphite (e.g. if running locally)",
		EnvVar: "GRAPHITE_TCP_ADDRESS",
	})
	graphitePrefix := app.String(cli.StringOpt{
		Name:   GraphitePrefix,
		Value:  "coco.services.$ENV.image-resolver.0",
		Desc:   "Prefix to use. Should start with content, include the environment, and the host name. e.g. coco.pre-prod.image-resolver.1",
		EnvVar: "GRAPHITE_PREFIX",
	})
	logMetrics := app.Bool(cli.BoolOpt{
		Name:   LogMetrics,
		Value:  false,
		Desc:   "Whether to log metrics. Set to true if running locally and you want metrics output",
		EnvVar: "LOG_METRICS",
	})
	embedsType := app.String(cli.StringOpt{
		Name:    EmbedsType,
		Value:  "http://www.ft.com/ontology/content/ImageSet",
		Desc:   "The type suported for embedded images, ex ImageSet",
		EnvVar: "EMBEDS_TYPE_WHITELIST",
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
			AppPort:             *port,
			Content_public_read: *cprHost,
			RouterAddress:       *routerAddress,
			GraphiteTCPAddress:  *graphiteTCPAddress,
			GraphitePrefix:      *graphitePrefix,
			HttpClient:          httpClient,
		}

		scMap := map[string]interface{}{
			Port:          port,
			CprHost:       cprHost,
			RouterAddress: routerAddress,
			GraphiteTCP:   graphiteTCPAddress,
			GraphitePrefix:             graphitePrefix,
		}

		var reader content.Reader
		var parser content.Parser
		var ir content.Resolver
		reader = content.NewContentReader(*cprHost, *routerAddress)
		parser = content.NewBodyParser(*embedsType)
		ir = content.NewImageResolver(&reader, &parser)

		baseftrwapp.OutputMetricsIfRequired(*graphiteTCPAddress, *graphitePrefix, *logMetrics)
		appLogger := content.NewAppLogger()
		ch := &content.ContentHandler{ServiceConfig:&sc, Service: ir, Log: appLogger}
		h := setupServiceHandler(sc, ch)
		appLogger.ServiceStartedEvent(AppCode, scMap)
		err := http.ListenAndServe(":"+*port, h)
		if err != nil {
			log.Fatalf("Unable to start server: %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func setupServiceHandler(sc content.ServiceConfig, ch *content.ContentHandler) *mux.Router {
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
