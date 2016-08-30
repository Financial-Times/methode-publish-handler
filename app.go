package main

import (
	"net/http"
	"os"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
	oldhttphandlers "github.com/Financial-Times/http-handlers-go/httphandlers"
	"github.com/Financial-Times/methode-publish-handler/vanities"
	"github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
)

const serviceDescription = "A RESTful API which accepts Methode Articles and appends a vanity url before passing on to the UPP Stack"

var timeout = time.Duration(10 * time.Second)

func main() {
	app := cli.App("methode-publish-handler", serviceDescription)
	app.Version("methode-publish-handler", "0.0.0")

	appName := app.StringOpt("app-name", "methode-publish-handler", "The name of this service")
	appPort := app.StringOpt("app-port", "8084", "HTTP Port for the app")

	notifierName := app.StringOpt("notifier", "cms-notifier", "The notifier service name")
	notifierURL := app.StringOpt("notifier-url", "http://localhost:8080/notify", "The url for the notifier")
	notifierPanicGuideURL := app.StringOpt("notifier-panic-guide", "", "The notifier panic guide url")
	notifierHealthcheckURL := app.StringOpt("notifier-health-url", "http://localhost:8080/__health", "The notifier healthcheck url")

	/*logMetrics := app.Bool(cli.BoolOpt{
		Name:   "log-metrics",
		Value:  false,
		Desc:   "Whether to log metrics. Set to true if running locally and you want metrics output",
		EnvVar: "LOG_METRICS",
	})*/ // TODO: check if needed

	app.Action = func() {
		sc := ServiceConfig{
			&http.Client{Timeout: timeout},
			*appName,
			*appPort,
			*notifierName,
			*notifierURL,
			*notifierHealthcheckURL,
			*notifierPanicGuideURL,
		}
		appLogger := NewAppLogger()
		metricsHandler := NewMetrics()
		notifierHandler := NotifierHandler{&sc, appLogger, &metricsHandler, &vanities.Vanity{}}

		handler := setupServiceHandler(sc, metricsHandler, notifierHandler)

		appLogger.ServiceStartedEvent(*appName, sc.asMap())
		//metricsHandler.OutputMetricsIfRequired(*graphiteTCPAddress, *graphitePrefix, *logMetrics) // TODO: check if needed

		err := http.ListenAndServe(":"+*appPort, handler)

		if err != nil {
			logrus.Fatalf("Unable to start server: %v", err)
		}
	}
	app.Run(os.Args)
}

func setupServiceHandler(config ServiceConfig, metricsHandler Metrics, notifierHandler NotifierHandler) *mux.Router {
	r := mux.NewRouter()

	decoratedHandler := oldhttphandlers.HTTPMetricsHandler(metricsHandler.registry, oldhttphandlers.TransactionAwareRequestLoggingHandler(logrus.StandardLogger(), notifierHandler))
	r.Path("/notify").Handler(handlers.MethodHandler{"POST": decoratedHandler})

	r.Path(httphandlers.BuildInfoPath).HandlerFunc(httphandlers.BuildInfoHandler)
	r.Path(httphandlers.PingPath).HandlerFunc(httphandlers.PingHandler)

	r.Path("/__health").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(fthealth.Handler(config.appName, serviceDescription, config.notifierCheck()))})
	r.Path("/__metrics").Handler(handlers.MethodHandler{"GET": http.HandlerFunc(metricsHttpEndpoint)})

	return r
}

type ServiceConfig struct {
	httpClient             *http.Client
	appName                string
	appPort                string
	notifierName           string
	notifierURL            string
	notifierHealthcheckURL string
	notifierPanicGuideURL  string
}

func (sc ServiceConfig) asMap() map[string]interface{} {
	return map[string]interface{}{
		"app-name":                 sc.appName,
		"app-port":                 sc.appPort,
		"notifier":                 sc.notifierName,
		"notifier-url":             sc.notifierURL,
		"notifier-health-url":      sc.notifierHealthcheckURL,
		"notifier-panic-guide-url": sc.notifierPanicGuideURL,
	}
}
