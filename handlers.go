package main

import (
	"encoding/json"
	"net/http"

	"io"

	"github.com/Financial-Times/methode-publish-handler/vanities"
	tid "github.com/Financial-Times/transactionid-utils-go"
	"golang.org/x/net/context"
)

const (
	uuidKey = "uuid"
)

// NotifierHandler Handles publish requests, and the vanity retrieval.
type NotifierHandler struct {
	config  *ServiceConfig
	log     *AppLogger
	metrics *Metrics
}

type publishedArticle struct {
	uuid string
}

func (h NotifierHandler) ServeHTTP(responseWriter http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	article := publishedArticle{}
	decoder.Decode(&article)

	h.log.TransactionStartedEvent(r.RequestURI, tid.GetTransactionIDFromRequest(r), article.uuid)

	vanity := vanities.GetVanity()
	article = appendVanityToContent(article, vanity)

	ctx := tid.TransactionAwareContext(context.Background(), r.Header.Get(tid.TransactionIDHeader))
	ctx = context.WithValue(ctx, uuidKey, article.uuid)

	ok, resp := h.postToNotifier(ctx, responseWriter)
	if !ok {
		return
	}

	io.Copy(responseWriter, resp.Body)
	h.metrics.recordResponseEvent()
}

func (h NotifierHandler) postToNotifier(ctx context.Context, responseWriter http.ResponseWriter) (ok bool, resp *http.Response) {
	uuid := ctx.Value(uuidKey).(string)
	transactionId, _ := tid.GetTransactionIDFromContext(ctx)

	h.log.RequestEvent(h.config.notifierName, h.config.notifierURL, transactionId, uuid)
	req, err := http.NewRequest("POST", h.config.notifierURL, nil)

	req.Header.Set(tid.TransactionIDHeader, transactionId)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)

	if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK) {
		responseWriter.WriteHeader(resp.StatusCode)
	} else {
		responseWriter.WriteHeader(http.StatusServiceUnavailable)
	}

	if err != nil {
		h.log.ErrorEvent(h.config.notifierName, h.config.notifierURL, transactionId, err, uuid)
		h.metrics.recordErrorEvent()
		return false, nil
	}

	if resp.StatusCode == http.StatusOK {
		h.log.ResponseEvent(h.config.notifierName, req.URL.String(), resp, uuid)
		return true, resp
	}

	h.log.RequestFailedEvent(h.config.notifierName, h.config.notifierURL, resp, uuid)
	h.metrics.recordRequestFailedEvent()
	return false, resp
}

func appendVanityToContent(article publishedArticle, vanity vanities.Vanity) publishedArticle {
	return article
}
