package jsonrpc

import (
	"context"
	"encoding/json"

	"github.com/Pantani/errors"
	"github.com/Pantani/request"
)

var (
	_requestID = int64(0)
)

const (
	// Version default JSON-RPC version.
	Version = "2.0"
)

type (
	// Request represents a client request wrapper.
	Request struct {
		request.Request
	}

	// RPCRequests represents a list of JSON-RPC requests.
	RPCRequests []*RPCRequest

	// RPCRequest represents the JSON-RPC request object.
	RPCRequest struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
		ID      int64       `json:"id,omitempty"`
	}

	// RPCResponse represents the JSON-RPC response object.
	RPCResponse struct {
		JSONRPC string      `json:"jsonrpc"`
		Error   *RPCError   `json:"error,omitempty"`
		Result  interface{} `json:"result,omitempty"`
		ID      int64       `json:"id,omitempty"`
	}

	// RPCError represents the JSON-RPC error object.
	RPCError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

// getObject bind the request result into the struct.
// It returns an error if occurs.
func (r *RPCResponse) getObject(toType interface{}) error {
	js, err := json.Marshal(r.Result)
	if err != nil {
		return errors.E(err, "json-rpc getObject Marshal error", errors.Params{"obj": toType})
	}

	err = json.Unmarshal(js, toType)
	if err != nil {
		return errors.E(err, "json-rpc getObject Unmarshal error", errors.Params{"obj": toType, "string": string(js)})
	}
	return nil
}

// RPCCall make a JSON-RPC request and bind the result into the generic interface.
// E.g.:
//	var tx []Tx
//	err = c.RpcCall(&tx, "getTransactionsByAddress", []string{address, "25"})
// It returns an error if occurs.
func (r *Request) RPCCall(result interface{}, method string, params interface{}) error {
	req := &RPCRequest{JSONRPC: Version, Method: method, Params: params, ID: genID()}
	var resp *RPCResponse
	err := r.Post(&resp, "", req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.E("RPC Call error", errors.Params{
			"method":        method,
			"error_code":    resp.Error.Code,
			"error_message": resp.Error.Message})
	}
	return resp.getObject(result)
}

// RPCCallWithContext make a JSON-RPC request and bind the result into the generic interface passing the context.
// E.g.:
//	var tx []Tx
//	err = c.RpcCall(context.Background(), &tx, "getTransactionsByAddress", []string{address, "25"})
// It returns an error if occurs.
func (r *Request) RPCCallWithContext(ctx context.Context, result interface{}, method string, params interface{}) error {
	req := &RPCRequest{JSONRPC: Version, Method: method, Params: params, ID: genID()}
	var resp *RPCResponse
	err := r.PostWithContext(ctx, &resp, "", req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return errors.E("RPC Call error", errors.Params{
			"method":        method,
			"error_code":    resp.Error.Code,
			"error_message": resp.Error.Message})
	}
	return resp.getObject(result)
}

// RPCBatchCall make a batch of JSON-RPC requests.
// E.g.:
//	var requests RpcRequests
//	for _, hash := range hashes {
//		requests = append(requests, &RpcRequest{
//			Method: "GetTransaction",
//			Params: []string{hash},
//		})
//	}
//	responses, err := c.RpcBatchCall(requests)
// It returns the result and an error if occurs.
func (r *Request) RPCBatchCall(requests RPCRequests) ([]RPCResponse, error) {
	var resp []RPCResponse
	err := r.Post(&resp, "", requests.fillDefaultValues())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RPCBatchCallWithContext make a batch of JSON-RPC requests with context.
// E.g.:
//	var requests RpcRequests
//	for _, hash := range hashes {
//		requests = append(requests, &RpcRequest{
//			Method: "GetTransaction",
//			Params: []string{hash},
//		})
//	}
//	responses, err := c.RpcBatchCall(context.Background(), requests)
// It returns the result and an error if occurs.
func (r *Request) RPCBatchCallWithContext(ctx context.Context, requests RPCRequests) ([]RPCResponse, error) {
	var resp []RPCResponse
	err := r.PostWithContext(ctx, &resp, "", requests.fillDefaultValues())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// fillDefaultValues fill the default JSON-RPC parameters.
func (rs RPCRequests) fillDefaultValues() RPCRequests {
	for _, r := range rs {
		r.JSONRPC = Version
		r.ID = genID()
	}
	return rs
}

// genID generate the request ID.
func genID() int64 {
	_requestID++
	return _requestID
}
