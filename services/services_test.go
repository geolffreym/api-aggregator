package services

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockProvider struct {
}

func (*MockProvider) TxByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request) {}
func (*MockProvider) TxByBlockHashAndIndex(w http.ResponseWriter, r *http.Request)   {}
func (*MockProvider) SendRawTransaction(w http.ResponseWriter, r *http.Request)      {}
func (*MockProvider) BlockByNumber(w http.ResponseWriter, r *http.Request)           {}
func (*MockProvider) BlockByHash(w http.ResponseWriter, r *http.Request)             {}
func (*MockProvider) BlockNumber(w http.ResponseWriter, r *http.Request)             {}

var mock *MockProvider

func getRequest(url string, t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Creating 'GET %s' request failed!", url)
	}

	return req
}

func setup() (*Services, *httptest.ResponseRecorder, *mux.Router) {
	m := mux.NewRouter()
	v1 := m.PathPrefix("/v1").Subrouter()

	mock = &MockProvider{}
	service := New(v1, mock)
	recorder := httptest.NewRecorder()
	return service, recorder, v1
}

func TestNewServiceWithValidRouter(t *testing.T) {
	service, _, v1 := setup()
	if service.r != v1 {
		t.Errorf("expected service router to be %v, got %v", v1, service)
	}
}

func TestNewServiceWithValidProvider(t *testing.T) {
	service, _, _ := setup()
	if service.p != mock {
		t.Errorf("expected service provider to be %v, got %v", mock, service)
	}
}

func TestDisabledSendTransactionEndpoint(t *testing.T) {
	service, recorder, v1 := setup()
	uri := "/v1/send/raw"
	// enable routes but not enable "send" service
	service.EnableBlocks()
	req, err := http.NewRequest("POST", uri, bytes.NewReader([]byte("test")))

	if err != nil {
		t.Fatalf("creating 'POST %s' request failed!", uri)
	}

	v1.ServeHTTP(recorder, req)
	// Expected 404 for not enabled service
	if recorder.Code != http.StatusNotFound {
		t.Error("server error: Returned ", recorder.Code, " instead of ", http.StatusNotFound)
	}
}

func TestBlockByNumberEndpoint(t *testing.T) {
	// Initialize provider RPC service
	service, recorder, v1 := setup()
	service.EnableBlocks()

	req := getRequest("/v1/block/by/number/0x5BAD55/true", t)
	v1.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Error("server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}

}

func TestBlockByHashEndpoint(t *testing.T) {
	// Initialize provider RPC service
	service, recorder, v1 := setup()
	service.EnableBlocks()

	req := getRequest("/v1/block/by/hash/0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35/true", t)
	v1.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Error("server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}
}

func TestTransactionByNumber(t *testing.T) {
	// Initialize provider RPC service
	service, recorder, v1 := setup()
	service.EnableTransactions()

	req := getRequest("/v1/tx/by/number/0x5BAD55/0x0", t)
	v1.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Error("server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}
}

func TestTransactionByHash(t *testing.T) {
	// Initialize provider RPC service
	service, recorder, v1 := setup()
	service.EnableTransactions()

	req := getRequest("/v1/tx/by/hash/0xb3b20624f8f0f86eb50dd04688409e5cea4bd02d700bf6e79e9384d47d6a5a35/0x0", t)
	v1.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Error("server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}
}
