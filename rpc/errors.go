// Application level errors for RPC

package rpc

import (
	"github.com/gorilla/rpc/v2/json2"
)

// ErrNodeFetch is an error to denote a generic error encountered when fetching a node
var ErrNodeFetch = &json2.Error{
	Code:    1000,
	Message: "Error fetching node",
}

// ErrNodeForksFetch denotes error when fetching forks of a node
var ErrNodeForksFetch = &json2.Error{
	Code:    1001,
	Message: "Error fetching forks of node",
}

// ErrNodesFetch denoted error encountered when fetching nodes of a thread
var ErrNodesFetch = &json2.Error{
	Code:    1003,
	Message: "Error fetching nodes of a thread",
}

// ErrCreateThread denotes error encountered when creating a thread
var ErrCreateThread = &json2.Error{
	Code:    1004,
	Message: "Error creating thread",
}

// ErrCreateNodeRevision denotes error encountered when creating revision for a node
var ErrCreateNodeRevision = &json2.Error{
	Code:    1005,
	Message: "Error creating node revision",
}

// ErrAddComment denotes error encountered when adding a comment
var ErrAddComment = &json2.Error{
	Code:    1006,
	Message: "Error adding comment",
}

// ErrNodeRevisionFetch denotes error encountered when fetching a node revision
var ErrNodeRevisionFetch = &json2.Error{
	Code:    1007,
	Message: "Error fetching node revision",
}

// ErrNodeRevisionsFetch denotes error encountered when fetching node revisions
var ErrNodeRevisionsFetch = &json2.Error{
	Code:    1008,
	Message: "Error fetching node revisions",
}

// ErrNodeRepliesFetch denotes error encountered when fetching node replies
var ErrNodeRepliesFetch = &json2.Error{
	Code:    1009,
	Message: "Error fetching node replies",
}

// ErrThreadForksFetch denotes error encountered when fetching forks in a thread
var ErrThreadForksFetch = &json2.Error{
	Code:    1010,
	Message: "Error fetching forks in a thread",
}

// ErrNodeFork denotes error encountered when forking node into thread
var ErrNodeFork = &json2.Error{
	Code:    1011,
	Message: "Error forking node into a thread",
}
