package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Financial-Times/methode-publish-handler/vanities"
	"github.com/stretchr/testify/assert"
)

var logger = NewAppLogger()
var testClient = &http.Client{Timeout: timeout}

type MockVanityService struct{}

func (MockVanityService) GetVanity() vanities.Vanity {
	return vanities.Vanity{"I am vain"}
}

func createTestHandler(mockExternalServer *httptest.Server, vanityService vanities.VanityService) NotifierHandler {
	m := NewMetrics()
	conf := ServiceConfig{testClient, "", "", "cms-notifier", mockExternalServer.URL, "", ""}
	if vanityService == nil {
		vanityService = vanities.Vanity{}
	}
	return NotifierHandler{&conf, logger, &m, vanityService}
}

func createMockExternalServer(assert *assert.Assertions, mockStatus int, mockBody string, expectedRequestBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.log.Info("Request received from handler. Mocking response.")
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(expectedRequestBody, string(body))
		w.WriteHeader(mockStatus)
		io.WriteString(w, mockBody)
	}))
}

func createMockServer(handler NotifierHandler) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler.ServeHTTP))
}

func TestHandlerProxiesRequestsToNotifierAndAddsVanity(t *testing.T) {
	assert := assert.New(t)
	expectedBody := "{}"
	expectedStatus := http.StatusOK

	mockRequestBody := `{"uuid":"6fcf1090-67b7-11e6-a0b1-d87a9fea034f","systemAttributes":"Some sys attrs","lastModified":"Modified at some stage recently","type":"An interesting document type","workflowStatus":"WorkflowStatus = Currently being unit tested","usageTickets":"Ticket number 42","linkedObjects":["One object linked through this string","Another object"],"value":"Some wonderful FT Content","attributes":"Some other attributes, which probably aren't system related"}`
	expectedRequestBody := `{"uuid":"6fcf1090-67b7-11e6-a0b1-d87a9fea034f","systemAttributes":"Some sys attrs","lastModified":"Modified at some stage recently","type":"An interesting document type","workflowStatus":"WorkflowStatus = Currently being unit tested","usageTickets":"Ticket number 42","linkedObjects":["One object linked through this string","Another object"],"value":"Some wonderful FT Content","attributes":"Some other attributes, which probably aren't system related","webUrl":"I am vain"}`

	mockExternalServer := createMockExternalServer(assert, expectedStatus, expectedBody, expectedRequestBody)
	defer mockExternalServer.Close()
	handler := createTestHandler(mockExternalServer, &MockVanityService{})

	mockServer := createMockServer(handler)
	defer mockServer.Close()

	req, err := http.NewRequest("POST", mockServer.URL+"/notify", strings.NewReader(mockRequestBody))

	result, err := testClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal("application/json", result.Header.Get("Content-Type"))
	assert.Equal(expectedBody, string(body))
	assert.Equal(expectedStatus, result.StatusCode)
}

func TestHandlerDoesNotAddVanity(t *testing.T) {
	assert := assert.New(t)
	expectedBody := "{}"
	expectedStatus := http.StatusOK
	mockRequestBody := `{"uuid":"6fcf1090-67b7-11e6-a0b1-d87a9fea034f","systemAttributes":"Some sys attrs","lastModified":"Modified at some stage recently","type":"An interesting document type","workflowStatus":"WorkflowStatus = Currently being unit tested","usageTickets":"Ticket number 42","linkedObjects":["One object linked through this string","Another object"],"value":"Some wonderful FT Content","attributes":"Some other attributes, which probably aren't system related","webUrl":""}`

	mockExternalServer := createMockExternalServer(assert, expectedStatus, expectedBody, mockRequestBody)
	defer mockExternalServer.Close()
	handler := createTestHandler(mockExternalServer, nil)

	mockServer := createMockServer(handler)
	defer mockServer.Close()

	req, err := http.NewRequest("POST", mockServer.URL+"/notify", strings.NewReader(mockRequestBody))

	result, err := testClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(expectedBody, string(body))
	assert.Equal(expectedStatus, result.StatusCode)
}

func TestNotFoundNotifierResponse(t *testing.T) {
	assert := assert.New(t)
	expectedStatus := http.StatusNotFound
	mockRequestBody := `{"uuid":"6fcf1090-67b7-11e6-a0b1-d87a9fea034f","systemAttributes":"Some sys attrs","lastModified":"Modified at some stage recently","type":"An interesting document type","workflowStatus":"WorkflowStatus = Currently being unit tested","usageTickets":"Ticket number 42","linkedObjects":["One object linked through this string","Another object"],"value":"Some wonderful FT Content","attributes":"Some other attributes, which probably aren't system related","webUrl":""}`

	mockExternalServer := createMockExternalServer(assert, expectedStatus, "{}", mockRequestBody)
	defer mockExternalServer.Close()

	handler := createTestHandler(mockExternalServer, nil)

	mockServer := createMockServer(handler)
	defer mockServer.Close()

	req, err := http.NewRequest("POST", mockServer.URL+"/notify", strings.NewReader(mockRequestBody))

	result, err := testClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal("", string(body))
	assert.Equal(expectedStatus, result.StatusCode)
}

func TestOtherFailedNotifierRequests(t *testing.T) {
	assert := assert.New(t)
	expectedStatus := http.StatusServiceUnavailable
	mockRequestBody := `{"uuid":"6fcf1090-67b7-11e6-a0b1-d87a9fea034f","systemAttributes":"Some sys attrs","lastModified":"Modified at some stage recently","type":"An interesting document type","workflowStatus":"WorkflowStatus = Currently being unit tested","usageTickets":"Ticket number 42","linkedObjects":["One object linked through this string","Another object"],"value":"Some wonderful FT Content","attributes":"Some other attributes, which probably aren't system related","webUrl":""}`

	mockExternalServer := createMockExternalServer(assert, 509, "{}", mockRequestBody)
	defer mockExternalServer.Close()

	handler := createTestHandler(mockExternalServer, nil)

	mockServer := createMockServer(handler)
	defer mockServer.Close()

	req, err := http.NewRequest("POST", mockServer.URL+"/notify", strings.NewReader(mockRequestBody))

	result, err := testClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal("", string(body))
	assert.Equal(expectedStatus, result.StatusCode)
}
