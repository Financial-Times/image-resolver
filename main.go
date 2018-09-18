package main

import (
	"net"
	"net/http"
	"os"
	"time"

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
	embeddedContentTypeWhitelist := app.String(cli.StringOpt{
		Name:   "embeddedContentTypeWhitelist",
		Value:  "^(http://www.ft.com/ontology/content/ImageSet)",
		Desc:   "The type supported for embedded images, ex ImageSet",
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
			ContentSourceURL:     *contentSourceHealthURL,
			HttpClient:           httpClient,
		}

		reader := content.NewContentReader(*contentSourceAppName, *contentSourceURL, httpClient)
		ir := content.NewImageResolver(reader, *embeddedContentTypeWhitelist, *apiHost)

		h := setupServiceHandler(ir, sc)
		err := http.ListenAndServe(":"+*port, h)
		if err != nil {
			log.Fatalf("Unable to start server: %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)
	log.Infof("Application started with args %s", os.Args)
	app.Run(os.Args)
}

func setupServiceHandler(s *content.ImageResolver, sc content.ServiceConfig) *mux.Router {
	r := mux.NewRouter()
	ch := &content.Handler{Service: s}
	r.HandleFunc("/content/image", ch.GetContentImages).Methods("POST")
	r.HandleFunc("/internalcontent/image", ch.GetLeadImages).Methods("POST")

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
