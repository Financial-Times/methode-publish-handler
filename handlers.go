package main

import "net/http"

const uuidKey = "uuid"

type ContentHandler struct {
	serviceConfig *ServiceConfig
	log           *AppLogger
	metrics       *Metrics
}

func (h ContentHandler) handleError(w http.ResponseWriter, err error, serviceName string, url string, transactionId string, uuid string) {
	w.WriteHeader(http.StatusServiceUnavailable)
	h.log.ErrorEvent(serviceName, url, transactionId, err, uuid)
	h.metrics.recordErrorEvent()
}

func (h ContentHandler) handleFailedRequest(w http.ResponseWriter, resp *http.Response, serviceName string, url string, uuid string) {
	w.WriteHeader(http.StatusServiceUnavailable)
	h.log.RequestFailedEvent(serviceName, url, resp, uuid)
	h.metrics.recordRequestFailedEvent()
}

func (h ContentHandler) handleNotFound(w http.ResponseWriter, resp *http.Response, serviceName string, url string, uuid string) {
	w.WriteHeader(http.StatusNotFound)
	h.log.RequestFailedEvent(serviceName, url, resp, uuid)
	h.metrics.recordRequestFailedEvent()
}
