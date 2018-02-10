// Application level errors for RPC

package rpc

import (
	"github.com/gorilla/rpc/v2/json2"
)

// ErrNodeFetch is an error to denote a generic error encountered when fetching a node.
var ErrNodeFetch = &json2.Error{
	Code:    100,
	Message: "Error fetching node",
}
