// The services package is based on decoupling the methods provided by the providers.
// The strategy is to separate the different methods that share common features and allow them to be enabled or disabled.
// Use this package to define the routes for your services.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods
package services

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Provider interface define method to query providers.
type Provider interface {
	TxByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request)
	TxByBlockHashAndIndex(w http.ResponseWriter, r *http.Request)
	SendRawTransaction(w http.ResponseWriter, r *http.Request)
	BlockByNumber(w http.ResponseWriter, r *http.Request)
	BlockByHash(w http.ResponseWriter, r *http.Request)
	BlockNumber(w http.ResponseWriter, r *http.Request)
}

type Services struct {
	p Provider // Provider interface
	r *mux.Router
}

// New creates a new Services instance.
func New(router *mux.Router, provider Provider) *Services {
	// Convention to define sub-router for each group type of common factor methods.
	// eg: eth_get[Block]ByNumber, eth_get[Block]ByHash

	// In this case of not sharing any group could be handled directly.
	// versioned.HandleFunc("/accounts", infura.Accounts).Methods("GET")
	// versioned.HandleFunc("/gasPrice", infura.GasPrice).Methods("GET")

	// Api version handling in URI
	return &Services{provider, router}

}

// EnableBlocks enable service for "block" group.
func (s *Services) EnableBlocks() *Services {
	log.Print("blocks service enabled")
	// Groups in this context eg. Uncle, Transaction, Block, Send.
	// Routes for block group. eg: eth_get[Block]ByNumber, eth_get[Block]ByHash
	block := s.r.PathPrefix("/block").Subrouter()
	block.HandleFunc("/", s.p.BlockNumber).Methods(http.MethodGet)

	// Subrouter to group "by" condition
	by := block.PathPrefix("/by").Subrouter()
	by.HandleFunc("/number/{block}/{flag}", s.p.BlockByNumber).Methods(http.MethodGet)
	by.HandleFunc("/hash/{block}/{flag}", s.p.BlockByHash).Methods(http.MethodGet)
	return s
}

// EnableTransactions enable service for "transaction" group.
func (s *Services) EnableTransactions() *Services {
	log.Print("transactions service enabled")
	// Groups in this context eg. Uncle, Transaction, Block, Send.
	// Routes for transactions group eg. eth_getTransactionByBlockNumberAndIndex
	transactions := s.r.PathPrefix("/tx").Subrouter()

	// Subrouter to group "by" condition
	by := transactions.PathPrefix("/by").Subrouter()
	by.HandleFunc("/number/{block}/{index}", s.p.TxByBlockNumberAndIndex).Methods(http.MethodGet)
	by.HandleFunc("/hash/{block}/{index}", s.p.TxByBlockHashAndIndex).Methods(http.MethodGet)
	return s
}

// EnableTSend enable service for "send" group.
func (s *Services) EnableSendTransactions() *Services {
	log.Print("send service enabled")
	// Groups in this context eg. Uncle, Transaction, Block, Send.
	// Routes for transactions group eg. eth_getTransactionByBlockNumberAndIndex
	transactions := s.r.PathPrefix("/send").Subrouter()
	transactions.HandleFunc("/raw", s.p.SendRawTransaction).Methods(http.MethodPost)
	return s
}
