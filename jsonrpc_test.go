package jsonrpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRpcRequests_fillDefaultValues(t *testing.T) {
	tests := []struct {
		name string
		rs   RpcRequests
		want RpcRequests
	}{
		{
			"test 1",
			RpcRequests{{Method: "method1", Params: "params1"}},
			RpcRequests{{Method: "method1", Params: "params1", JsonRpc: JsonRpcVersion, Id: "method1"}},
		}, {
			"test 2",
			RpcRequests{
				{Method: "method1", Params: "params1"}, {Method: "method2", Params: "params2"}},
			RpcRequests{
				{Method: "method1", Params: "params1", JsonRpc: JsonRpcVersion, Id: "method1"},
				{Method: "method2", Params: "params2", JsonRpc: JsonRpcVersion, Id: "method2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rs.fillDefaultValues()
			assert.Equal(t, tt.want, got)
		})
	}
}
