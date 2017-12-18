package rpc

import (
	"log"
	"net/http"

	"github.com/applait/upspeak/models"
)

type NodeService struct{}

//////////////////////////
// Node reply construct //
//////////////////////////
type NodeReply struct {
}

////////////////////////////
// Node revision constuct //
////////////////////////////
type NodeRevisionReply struct {
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	CommitTime  string `json:"commit_time"`
	CommitterID int64  `json:"committer_id"`
}

//////////////
// Get node //
//////////////
type NodeGetArgs struct {
	NodeID int64 `json:"node_id"`
}

// Get gets a node
func (n *NodeService) Get(r *http.Request, args *NodeGetArgs, reply *models.GetNodeSchema) error {
	node := new(models.GetNodeSchema)
	node.NodeId = args.NodeID
	node, err := models.GetNode(node)
	if err != nil {
		log.Println(err)
	}
	return nil
}

///////////////////////
// Get node revision //
///////////////////////
type NodeGetRevisionArgs struct {
	NodeId     int64 `json:"node_id"`
	RevisionId int64 `json:"revision_id"`
}

// GetRevision gets a node revision
func (n *NodeService) GetRevision(r *http.Request, args *NodeGetRevisionArgs, reply *NodeRevisionReply) error {
	// TODO: Fill up NodeGetRevisionReply
	return nil

}

////////////////////////
// Get node revisions //
////////////////////////
type NodeGetRevisionsArgs struct {
	NodeID int `json:"node_id"`
}

// GetRevisions gets last 10 node revisions unless specified otherwise
func (n *NodeService) GetRevisions(r *http.Request, args *NodeGetRevisionsArgs, reply *[]NodeRevisionReply) error {
	return nil
}

//////////////////////
// Get node replies //
//////////////////////
type NodeGetRepliesArgs struct {
	NodeID int `json:"node_id"`
}

// GetReplies gets last 10 replies unless specified otherwise
func (n *NodeService) GetReplies(r *http.Request, args *NodeGetRepliesArgs, reply *[]NodeReply) error {
	return nil
}

////////////////
// Edit node  //
////////////////
type NodeEditArgs struct {
	NodeId   int64 `json:"node_id"`
	AuthorId int64 `json:"author_id"`
	Content  int64 `json:"content"`
}

// Edit creates a node revision
func (n *NodeService) Edit(r *http.Request, args *NodeEditArgs, reply *NodeReply) error {
	return nil
}

///////////////
// Fork node //
///////////////
type NodeForkArgs struct {
	NodeId   int64 `json:"node_id"`
	AuthorId int64 `json:"author_id"`
	TeamId   int64 `json:"team_id"`
}

// Fork creates a new thread from current node
func (n *NodeService) Fork(r *http.Request, args *NodeEditArgs, reply *NodeReply) error {
	// s/NodeReply/ThreadReply
	return nil
}

////////////////////
// Create comment //
////////////////////
type NodeCreateCommentArgs struct {
	ThreadId int64 `json:"thread_id"`
	AuthorId int64 `json:"author_id"`
	TeamId   int64 `json:"team_id"`
}

// CreateComment creates a new comment
func (n *NodeService) CreateComment(r *http.Request, args *NodeCreateCommentArgs, reply *NodeReply) error {
	return nil
}

////////////////
// Node utils //
////////////////

// resolvePermissions returns node's permission matrix
func (n *NodeService) resolvePermissions(nodeId int64) {
	// TODO: Define permission matrix data structure
}
