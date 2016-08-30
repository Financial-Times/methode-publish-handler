package main

import (
	"fmt"
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
)

func (config *ServiceConfig) notifierCheck() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "No Articles from PortalPub will be published!",
		Name:             config.notifierName + " Availabililty Check",
		PanicGuide:       config.notifierPanicGuideURL,
		Severity:         1,
		TechnicalSummary: "Checks that \"" + config.notifierName + "\" Service is reachable. MOPH publishes articles to \"" + config.notifierName + "\" after pre-processing them for vanity urls.",
		Checker: func() (string, error) {
			return checkServiceAvailability(config.httpClient, config.notifierName, config.notifierHealthcheckURL, "", "")
		},
	}
}

func checkServiceAvailability(client *http.Client, serviceName string, healthUri string, auth string, hostHeader string) (string, error) {
	req, err := http.NewRequest("GET", healthUri, nil)
	if auth != "" {
		req.Header.Set("Authorization", "Basic "+auth)
	}

	if hostHeader != "" {
		req.Host = hostHeader
	}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s service is unreachable", serviceName), fmt.Errorf("%s service is unreachable", serviceName)
	}
	return "Ok", nil
}
