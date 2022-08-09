package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"aggregator/services"

	"github.com/gorilla/mux"
)

type MockProvider struct {
}

func (*MockProvider) BlockByHash(w http.ResponseWriter, r *http.Request)             {}
func (*MockProvider) BlockNumber(w http.ResponseWriter, r *http.Request)             {}
func (*MockProvider) TxByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request) {}
func (*MockProvider) TxByBlockHashAndIndex(w http.ResponseWriter, r *http.Request)   {}
func (*MockProvider) SendRawTransaction(w http.ResponseWriter, r *http.Request)      {}
func (*MockProvider) BlockByNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Write([]byte(params["block"] + "." + params["flag"]))
}

func runTransactionTest(v1 *mux.Router) *httptest.ResponseRecorder {
	// First lets check that the block endpoints works fine
	req, _ := http.NewRequest("GET", "/v1/tx/by/number/0x5BAD55/0x0", nil)
	recorder := httptest.NewRecorder()
	v1.ServeHTTP(recorder, req)
	return recorder
}

func TestIntegrationEnableDisableService(t *testing.T) {
	m := mux.NewRouter()
	v1 := m.PathPrefix("/v1").Subrouter()
	service := services.New(v1, &MockProvider{})
	// Lets try to request a non active service
	recorder := runTransactionTest(v1)
	if recorder.Code == http.StatusOK {
		t.Fatalf("unexpected status code for disabled tx services endpoints: %d", recorder.Code)
	}

	// Now lets active the tx service
	service.EnableTransactions()
	recorder = runTransactionTest(v1)
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status code for enabled tx services endpoints: %d", recorder.Code)
	}

}

func TestIntegrationEndpointResponse(t *testing.T) {
	m := mux.NewRouter()
	v1 := m.PathPrefix("/v1").Subrouter()
	service := services.New(v1, &MockProvider{})

	service.EnableBlocks()
	expected := "0x5BAD55.true"
	// Request a enabled service and try to check expected response
	req, _ := http.NewRequest("GET", "/v1/block/by/number/0x5BAD55/true", nil)
	recorder := httptest.NewRecorder()
	v1.ServeHTTP(recorder, req)

	if recorder.Body.String() != expected {
		t.Errorf("unexpected body: %s; expected: %s", recorder.Body.String(), expected)
	}

}
