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
func (routes *Services) EnableBlocks() *Services {
	log.Print("blocks service enabled")
	// Groups in this context eg. Uncle, Transaction, Block, Send.
	// Routes for block group. eg: eth_get[Block]ByNumber, eth_get[Block]ByHash
	block := routes.r.PathPrefix("/block").Subrouter()
	block.HandleFunc("/", routes.p.BlockNumber).Methods(http.MethodGet)

	// Subrouter to group "by" condition
	by := block.PathPrefix("/by").Subrouter()
	by.HandleFunc("/number/{block}/{flag}", routes.p.BlockByNumber).Methods(http.MethodGet)
	by.HandleFunc("/hash/{block}/{flag}", routes.p.BlockByHash).Methods(http.MethodGet)
	return routes
}

// EnableTransactions enable service for "transaction" group.
func (routes *Services) EnableTransactions() *Services {
	log.Print("transactions service enabled")
	// Groups in this context eg. Uncle, Transaction, Block, Send.
	// Routes for transactions group eg. eth_getTransactionByBlockNumberAndIndex
	transactions := routes.r.PathPrefix("/tx").Subrouter()

	// Subrouter to group "by" condition
	by := transactions.PathPrefix("/by").Subrouter()
	by.HandleFunc("/number/{block}/{index}", routes.p.TxByBlockNumberAndIndex).Methods(http.MethodGet)
	by.HandleFunc("/hash/{block}/{index}", routes.p.TxByBlockHashAndIndex).Methods(http.MethodGet)
	return routes
}

// EnableTSend enable service for "send" group.
func (routes *Services) EnableSendTransactions() *Services {
	log.Print("send service enabled")
	// Groups in this context eg. Uncle, Transaction, Block, Send.
	// Routes for transactions group eg. eth_getTransactionByBlockNumberAndIndex
	transactions := routes.r.PathPrefix("/send").Subrouter()
	transactions.HandleFunc("/raw", routes.p.SendRawTransaction).Methods(http.MethodPost)
	return routes
}
