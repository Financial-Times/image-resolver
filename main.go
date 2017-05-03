package main

import (
	"net/http"
	"github.com/jawher/mow.cli"
	"os"
	"github.com/gorilla/mux"
	"github.com/Financial-Times/image-resolver/content"
	log "github.com/Sirupsen/logrus"
	"github.com/Financial-Times/service-status-go/gtg"
	"github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/handlers"
	"net"
	"time"
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/Sirupsen/logrus"
	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
)


func main() {
	app := cli.App("image-resolver", "Image resolver - unroll images for a given article UUID")
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "9090",
		Desc:   "application port",
		EnvVar: "PORT"})
	cprHost := app.String(cli.StringOpt{
		Name:   "cprHost",
		Value:  "content-public-read",
		Desc:   "Content read app",
		EnvVar: "CONTENT_READ_ENDPOINT",
	})
	routerAddress := app.String(cli.StringOpt{
		Name:   "routerAddress",
		Value:  "localhost:8080",
		Desc:   "Vulcan host",
		EnvVar: "ROUTER_ADDRESS",
	})
	graphiteTCPAddress := app.String(cli.StringOpt{
		Name:   "graphite-tcp-address",
		Value:  "",
		Desc:   "Graphite TCP address, e.g. graphite.ft.com:2003. Leave as default if you do NOT want to output to graphite (e.g. if running locally)",
		EnvVar: "GRAPHITE_TCP_ADDRESS",
	})
	graphitePrefix := app.String(cli.StringOpt{
		Name:   "graphite-prefix",
		Value:  "coco.services.$ENV.content-preview.0",
		Desc:   "Prefix to use. Should start with content, include the environment, and the host name. e.g. coco.pre-prod.sections-rw-neo4j.1",
		EnvVar: "GRAPHITE_PREFIX",
	})
	logMetrics := app.Bool(cli.BoolOpt{
		Name:   "log-metrics",
		Value:  false,
		Desc:   "Whether to log metrics. Set to true if running locally and you want metrics output",
		EnvVar: "LOG_METRICS",
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
			*port,
			*cprHost,
			*routerAddress,
			*graphiteTCPAddress,
			*graphitePrefix,
			httpClient,
		}

		var reader content.Reader
		var parser content.Parser
		var ir content.Resolver
		reader = content.NewContentReader(*cprHost, *routerAddress)
		parser = content.BodyParser{}
		ir = content.NewImageResolver(&reader, &parser)

		baseftrwapp.OutputMetricsIfRequired(*graphiteTCPAddress, *graphitePrefix, *logMetrics)
		ch := &content.ContentHandler{&sc, ir}
		h := setupServiceHandler(sc, ch)
		err := http.ListenAndServe(":" + *port, h)
		if err != nil {
			logrus.Fatalf("Unable to start server: %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func setupServiceHandler(sc content.ServiceConfig, ch *content.ContentHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/content", ch.GetContentImages).Methods("POST")

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	checks := []fthealth.Check{sc.ContentCheck()}
	hc := fthealth.HealthCheck{SystemCode: "image-resolver", Name: "Image Resolver", Description: "Image Resolver", Checks: checks }
	r.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(fthealth.Handler(&hc))})
	gtgHandler := httphandlers.NewGoodToGoHandler(gtg.StatusChecker(sc.GtgCheck))
	r.Path("/__gtg").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(gtgHandler)})
	return r
}
