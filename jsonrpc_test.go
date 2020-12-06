package jsonrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCRequests_fillDefaultValues(t *testing.T) {
	tests := []struct {
		name string
		rs   RPCRequests
		want RPCRequests
	}{
		{
			"test 1",
			RPCRequests{{Method: "method1", Params: "params1"}},
			RPCRequests{{Method: "method1", Params: "params1", JSONRPC: Version, ID: 1}},
		}, {
			"test 2",
			RPCRequests{
				{Method: "method1", Params: "params1"}, {Method: "method2", Params: "params2"}},
			RPCRequests{
				{Method: "method1", Params: "params1", JSONRPC: Version, ID: 2},
				{Method: "method2", Params: "params2", JSONRPC: Version, ID: 3},
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
