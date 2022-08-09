// Package infura implements all the methods needed to satisfied Provider interface.
// Methods are corresponding to RPC API service provided by Infura.
package infura

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
)

// Allowed RPC methods are defined here.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods
const (
	BlockNumber             = "eth_blockNumber"
	BlockByHash             = "eth_getBlockByHash"
	BlockByNumber           = "eth_getBlockByNumber"
	SendRawTransaction      = "eth_sendRawTransaction"
	TxByBlockNumberAndIndex = "eth_getTransactionByBlockNumberAndIndex"
	TxByBlockHashAndIndex   = "eth_getTransactionByBlockHashAndIndex"
)

type InfuraRPC struct {
	client *rpc.Client
}

// Call execute rpc call to infura RPC Api.
// Return response received if successful otherwise return nil data and error.
func (r *InfuraRPC) Call(method string, args ...any) (any, error) {
	// container struct
	var container any
	// Call Infura API RPC method
	if err := r.client.Call(&container, method, args...); err != nil {
		// Zeroed indexBlock returned if error
		log.Printf("error calling RPC method %s: %v", method, err)
		return nil, err
	}

	return container, nil
}

// JSONFromCall execute a rpc call and return JSON unmarshal response.
func (r *InfuraRPC) JSONFromCall(method string, args ...any) ([]byte, error) {
	// Call Infura API RPC method
	raw, err := r.Call(method, args...)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(raw)
	if err != nil {
		log.Printf("error trying to unmarshal response: %v", err)
		return nil, err

	}

	return data, nil
}

// Infura implement RPC methods provided by Infura RPC Api docs.
type Infura struct {
	rpc *InfuraRPC
}

// parseFlag serve as helper to handle boolean params.
// if not boolean method is passed return default flag=true
func parseFlag(flag string) bool {
	// ref: https://go.dev/src/strconv/atob.go?s=351:391#L1
	f, err := strconv.ParseBool(flag)
	if err != nil {
		// default flag
		return true
	}

	return f
}

// New create a new Infura instance.
func New(endpoint string) (*Infura, error) {
	// Reuse client for Infura calls
	client, err := rpc.Dial(endpoint)
	if err != nil {
		log.Fatalf("could not connect to Infura: %v", err)
		return nil, err
	}

	rpc := &InfuraRPC{client}
	return &Infura{rpc}, nil
}

// TxByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionbyblocknumberandindex
func (i *Infura) TxByBlockNumberAndIndex(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	block := params["block"] // hexadecimal block number, or the string "latest", "earliest" or "pending".
	index := params["index"] // a hex of the integer representing the position in the block

	// Call Infura API eth_getTransactionByBlockNumberAndIndex RPC method
	response, err := i.rpc.JSONFromCall(TxByBlockNumberAndIndex, block, index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(response)
}

// TxByBlockHashAndIndex returns information about a transaction by block hash and transaction index position.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionbyblockhashandindex
func (i *Infura) TxByBlockHashAndIndex(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	block := params["block"] // a string representing the hash (32 bytes) of a block.
	index := params["index"] // a hex of the integer representing the position in the block

	// Call Infura API eth_getTransactionByBlockNumberAndIndex RPC method
	response, err := i.rpc.JSONFromCall(TxByBlockHashAndIndex, block, index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(response)
}

// BlockByNumber returns information about a block by hash.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_getblockbynumber
func (i *Infura) BlockByNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	number := params["block"]         // hexadecimal block number, or the string "latest", "earliest" or "pending".
	flag := parseFlag(params["flag"]) //  if set to true, it returns the full transaction objects, if false only the hashes of the transactions.

	// Call Infura API eth_getBlockByNumber RPC method
	response, err := i.rpc.JSONFromCall(BlockByNumber, number, flag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(response)

}

// BlockByHash returns information about a block by hash.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_getblockbyhash
func (i *Infura) BlockByHash(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	number := params["block"]         // hexadecimal block number, or the string "latest", "earliest" or "pending".
	flag := parseFlag(params["flag"]) // if set to true, it returns the full transaction objects, if false only the hashes of the transactions.

	// Call Infura API eth_getBlockByHash RPC method
	response, err := i.rpc.JSONFromCall(BlockByHash, number, flag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(response)

}

// BlockNumber returns the current "latest" block number.
// ref: https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_blocknumber
func (i *Infura) BlockNumber(w http.ResponseWriter, r *http.Request) {
	// Call Infura API eth_getBlockNumber RPC method
	response, err := i.rpc.JSONFromCall(BlockNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(response)

}

// SendRawTransaction submits a pre-signed transaction for broadcast to the Ethereum network.
func (i *Infura) SendRawTransaction(w http.ResponseWriter, r *http.Request) {
	tx := r.FormValue("tx") // The signed transaction data.

	// Call Infura API eth_sendRawTransaction RPC method
	response, err := i.rpc.JSONFromCall(SendRawTransaction, tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(response)

}
