# JSON-RPC request abstraction

Simple abstraction for JSON-RPC requests.

Initialize the client:
```go
import "github.com/Pantani/jsonrpc"

client := jsonrpc.InitRpcClient("http://127.0.0.1:8080")
```

## Methods

### Rpc Call

```go
var txs []Transaction
err := client.RpcCall(&txs, "getTransactionsByAddress", []string{"d48182276127b149a9710e78c436fb4bc1c4dc0b", "25"})
```

### Rpc Batch Call

```go
var requests jsonrpc.RpcRequests
for _, hash := range hashes {
    requests = append(requests, &jsonrpc.RpcRequest{
        Method: "GetTransaction",
        Params: []string{hash},
    })
}
responses, err := client.RpcBatchCall(requests)
if err != nil {
    panic(err)
}
for _, result := range responses {
    var tx Transaction
    if mapstructure.Decode(result.Result, &tx) != nil {
        continue
    }
    txs = append(txs, tx)
}
```
