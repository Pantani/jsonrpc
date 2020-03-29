package jsonrpc

import (
	"encoding/json"
	"github.com/Pantani/errors"
	"github.com/Pantani/request"
)

type Request struct {
	request.Request
}

const (
	JsonRpcVersion = "2.0"
)

type (
	RpcRequests []*RpcRequest
	RpcRequest  struct {
		JsonRpc string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
		Id      string      `json:"id,omitempty"`
	}

	RpcResponse struct {
		JsonRpc string      `json:"jsonrpc"`
		Error   *RpcError   `json:"error,omitempty"`
		Result  interface{} `json:"result,omitempty"`
		Id      string      `json:"id,omitempty"`
	}

	RpcError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

// InitRpcClient initialize the rpc client.
// It returns the rpc request object.
func InitRpcClient(baseUrl string) Request {
	return Request{
		request.Request{
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			HttpClient:   request.DefaultClient,
			ErrorHandler: request.DefaultErrorHandler,
			BaseUrl:      baseUrl,
		},
	}
}

// GetObject parse the result object into a generic interface.
// It returns an error if occurs.
func (r *RpcResponse) GetObject(toType interface{}) error {
	js, err := json.Marshal(r.Result)
	if err != nil {
		return errors.E(err, "json-rpc GetObject Marshal error", errors.Params{"obj": toType})
	}

	err = json.Unmarshal(js, toType)
	if err != nil {
		return errors.E(err, "json-rpc GetObject Unmarshal error", errors.Params{"obj": toType, "string": string(js)})
	}
	return nil
}

// RpcCall make a rpc call and return the result inside a generic interface bind.
// It returns an error if occurs.
func (r *Request) RpcCall(result interface{}, method string, params interface{}) error {
	req := &RpcRequest{JsonRpc: JsonRpcVersion, Method: method, Params: params, Id: method}
	var resp *RpcResponse
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
	return resp.GetObject(result)
}

// RpcBatchCall make a rpc batch call and return the result list.
// It returns a result list and an error if occurs.
func (r *Request) RpcBatchCall(requests RpcRequests) ([]RpcResponse, error) {
	var resp []RpcResponse
	err := r.Post(&resp, "", requests.fillDefaultValues())
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// fillDefaultValues fill the default values for a RPC request.
func (rs RpcRequests) fillDefaultValues() RpcRequests {
	for _, r := range rs {
		r.JsonRpc = JsonRpcVersion
		r.Id = r.Method
	}
	return rs
}
