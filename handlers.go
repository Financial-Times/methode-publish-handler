package main

import (
	"bytes"
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
	config        *ServiceConfig
	log           *AppLogger
	metrics       *Metrics
	vanityService vanities.VanityService
}

type PublishedArticle struct {
	UUID             string   `json:"uuid,omitempty"`
	SystemAttributes string   `json:"systemAttributes,omitempty"`
	LastModified     string   `json:"lastModified,omitempty"`
	Type             string   `json:"type,omitempty"`
	WorkflowStatus   string   `json:"workflowStatus,omitempty"`
	UsageTickets     string   `json:"usageTickets,omitempty"`
	LinkedObjects    []string `json:"linkedObjects,omitempty"`
	Value            string   `json:"value,omitempty"`
	Attributes       string   `json:"attributes,omitempty"`
	WebURL           string   `json:"webUrl,omitempty"`
}

func (h NotifierHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	article := PublishedArticle{}
	decoder.Decode(&article)

	h.log.TransactionStartedEvent(r.RequestURI, tid.GetTransactionIDFromRequest(r), article.UUID)

	vanity := h.vanityService.GetVanity()
	article = appendVanityToContent(article, vanity)

	ctx := tid.TransactionAwareContext(context.Background(), r.Header.Get(tid.TransactionIDHeader))
	ctx = context.WithValue(ctx, uuidKey, article.UUID)

	ok, resp := h.postToNotifier(ctx, w, article)
	if !ok {
		return
	}

	io.Copy(w, resp.Body)
	h.metrics.recordResponseEvent()
}

func (h NotifierHandler) postToNotifier(ctx context.Context, w http.ResponseWriter, article PublishedArticle) (ok bool, resp *http.Response) {
	uuid := ctx.Value(uuidKey).(string)
	transactionID, _ := tid.GetTransactionIDFromContext(ctx)

	js, err := json.Marshal(article)
	if err != nil {
		h.log.ErrorEvent(h.config.notifierName, h.config.notifierURL, transactionID, err, uuid)
		h.metrics.recordErrorEvent()
		return false, nil
	}

	h.log.RequestEvent(h.config.notifierName, h.config.notifierURL, transactionID, uuid)
	req, err := http.NewRequest("POST", h.config.notifierURL, bytes.NewBuffer(js))

	req.Header.Set(tid.TransactionIDHeader, transactionID)
	req.Header.Set("Content-Type", "application/json")

	resp, err = h.config.httpClient.Do(req)

	if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusOK) {
		w.WriteHeader(resp.StatusCode)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err != nil {
		h.log.ErrorEvent(h.config.notifierName, h.config.notifierURL, transactionID, err, uuid)
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

func appendVanityToContent(article PublishedArticle, vanity vanities.Vanity) PublishedArticle {
	article.WebURL = vanity.WebURL
	return article
}
