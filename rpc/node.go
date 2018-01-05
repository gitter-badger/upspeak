package rpc

import (
	"log"
	"net/http"

	"github.com/applait/upspeak/models"
)

type NodeService struct{}

//////////////
// Get node //
//////////////
type NodeGetArgs struct {
	NodeID int64 `json:"node_id"`
}

type NodeGetReply struct {
	Node *models.Node `json:"node"`
}

// Get gets a node
func (n *NodeService) Get(r *http.Request, args *NodeGetArgs, reply *NodeGetReply) error {
	_node, err := models.GetNode(&args.NodeID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.Node = _node
	return nil
}

///////////////////////
// Get forks of node //
///////////////////////
type NodeGetForksOfNodeArgs struct {
	NodeID int64 `json:"node_id"`
}

type NodeGetForksOfNodeReply []struct {
	NodeID int64 `json:"node_id"`
}

func (n *NodeService) GetForksOfNode(r *http.Request, args *NodeGetForksOfNodeArgs, reply *NodeGetForksOfNodeReply) error {
	nodes, err := models.GetForksOfNode(args.NodeID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply = new(NodeGetForksOfNodeReply)
	for i, n := range nodes {
		(*reply)[i].NodeID = n.NodeID
	}

	return nil
}

//////////////////////////
// List nodes of thread //
//////////////////////////

type ThreadNodesArgs struct {
	ThreadID int64 `json:"thread_id"`
}

type ThreadNodesReply struct {
	ThreadID int64          `json:"thread_id"`
	Count    int            `json:"count"`
	Nodes    []*models.Node `json:"nodes"`
}

// ThreadNodes lists nodes for a given thread
func (n *NodeService) ThreadNodes(r *http.Request, args *ThreadNodesArgs, reply *ThreadNodesReply) error {
	nodes, err := models.GetNodes(&args.ThreadID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.ThreadID = args.ThreadID
	reply.Count = len(nodes)
	reply.Nodes = nodes
	return nil
}

///////////////////
// Create thread //
///////////////////
type CreateThreadArgs struct {
	TeamID   int64             `json:"team_id"`
	UserID   int64             `json:"user_id"`
	DataType string            `json:"data_type"`
	Subject  models.NullString `json:"subject"`
	Body     models.NullString `json:"body"`
	RichData models.JSONB      `json:"rich_data"`
}

type CreateThreadReply struct {
	ThreadID int64 `json:"thread_id"`
}

func (n *NodeService) CreateThread(r *http.Request, args *CreateThreadArgs, reply *CreateThreadReply) error {
	threadID, err := models.CreateThread(&models.CreateThreadSchema{
		TeamID: args.TeamID,
		UserID: args.UserID,
		Data: models.NodeData{
			DataType: args.DataType,
			Subject:  args.Subject,
			Body:     args.Body,
			RichData: args.RichData,
		},
	})
	if err != nil {
		log.Println(err)
	}

	reply.ThreadID = threadID
	return nil

}

////////////////
// Edit node  //
////////////////
type NodeEditArgs struct {
	NodeID  int64  `json:"node_id"`
	Subject string `json:"subject"`
	UserID  int64  `json:"user_id"`
	Body    string `json:"body"`
}

type NodeEditReply struct {
	NodeID int64 `json:"node_id"`
}

// Edit creates a node revision
func (n *NodeService) Edit(r *http.Request, args *NodeEditArgs, reply *NodeEditReply) error {
	revision := &models.AddNodeRevisionSchema{
		NodeID:  args.NodeID,
		Subject: args.Subject,
		UserID:  args.UserID,
		Body:    args.Body,
	}

	revision, err := models.AddNodeRevision(revision)
	if err != nil {
		log.Println(revision)
	}

	reply.NodeID = revision.NodeID
	return nil
}

/////////////////
// Add comment //
/////////////////
type NodeAddCommentArgs struct {
	ThreadID    int64        `json:"thread_id"`
	InReplyToID int64        `json:"in_reply_to_id"`
	DataType    string       `json:"data_type"`
	Subject     string       `json:"subject"`
	Body        string       `json:"body"`
	Extra       models.JSONB `json:"extra"`
	UserID      int64        `json:"user_id"`
}

type NodeAddCommentReply struct {
	NodeID int64 `json:"node_id"`
}

// AddComment adds a comment
func (n *NodeService) AddComment(r *http.Request, args *NodeAddCommentArgs, reply *NodeAddCommentReply) error {
	revision := &models.CreateCommentSchema{
		ThreadID:    args.ThreadID,
		InReplyToID: args.InReplyToID,
		DataType:    args.DataType,
		Subject:     args.Subject,
		Body:        args.Body,
		UserID:      args.UserID,
		Extra:       args.Extra,
	}

	revision, err := models.CreateComment(revision)
	if err != nil {
		log.Println(revision)
	}

	reply.NodeID = revision.NodeID
	return nil
}

///////////////////////
// Get node revision //
///////////////////////
type NodeGetRevisionArgs struct {
	NodeID int64 `json:"node_id"`
}

type NodeGetRevisionReply struct {
	NodeID int64 `json:"node_id"`
}

// GetRevision gets a node revision
func (n *NodeService) GetRevision(r *http.Request, args *NodeGetRevisionArgs, reply *NodeGetRevisionReply) error {
	node := &models.GetNodeRevisionSchema{
		NodeID: args.NodeID,
	}
	node, err := models.GetNodeRevision(node)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.NodeID = node.NodeID
	return nil
}

////////////////////////
// Get node revisions //
////////////////////////
type NodeGetRevisionsArgs struct {
	NodeID int64 `json:"node_id"`
}

type NodeGetRevisionsReply []struct {
	NodeID int64 `json:"node_id"`
}

// GetRevisions gets revisions of a node
func (n *NodeService) GetRevisions(r *http.Request, args *NodeGetRevisionsArgs, reply *NodeGetRevisionsReply) error {
	revisions, err := models.GetNodeRevisions(args.NodeID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply = new(NodeGetRevisionsReply)
	for i, n := range revisions {
		(*reply)[i].NodeID = n.NodeID
	}

	return nil

}

//////////////////////
// Get node replies //
//////////////////////
type NodeGetRepliesArgs struct {
	ThreadID int64 `json:"thread_id"`
}

type NodeGetRepliesReply []struct {
	NodeID int64 `json:"node_id"`
}

// GetReplies gets replies of a node
func (n *NodeService) GetReplies(r *http.Request, args *NodeGetRepliesArgs, reply *NodeGetRepliesReply) error {
	nodeReplies, err := models.GetReplies(args.ThreadID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply = new(NodeGetRepliesReply)
	for i, n := range nodeReplies {
		(*reply)[i].NodeID = n.NodeID
	}

	return nil

}

///////////////////////////
// Get forks in a thread //
///////////////////////////
type NodeGetForksInAThreadArgs struct {
	ThreadID int64 `json:"thread_id"`
}

type NodeGetForksInAThreadReply []struct {
	NodeID int64 `json:"node_id"`
}

func (n *NodeService) GetForksInAThread(r *http.Request, args *NodeGetForksInAThreadArgs, reply *NodeGetForksInAThreadReply) error {
	forks, err := models.GetForksInAThread(args.ThreadID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply = new(NodeGetForksInAThreadReply)
	for i, n := range forks {
		(*reply)[i].NodeID = n.NodeID
	}

	return nil
}

///////////////
// Fork node //
///////////////
type NodeForkNodeArgs struct {
	TeamID     int64  `json:"team_id"`
	UserID     int64  `json:"user_id"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
	ForkedFrom int64  `json:"forked_from"`
}

type NodeForkNodeReply struct {
	NodeID int64 `json:"node_id"`
}

// ForkNode creates a new thread from current node
func (n *NodeService) ForkNode(r *http.Request, args *NodeForkNodeArgs, reply *NodeForkNodeReply) error {
	thread := new(models.ForkNodeSchema)

	thread, err := models.ForkNode(args.TeamID, args.UserID, args.Subject, args.Body, args.ForkedFrom)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.NodeID = thread.NodeID
	return nil
}

////////////////
// Node utils //
////////////////

// resolvePermissions returns node's permission matrix
func (n *NodeService) resolvePermissions(nodeID int64) {
	// TODO: Define permission matrix data structure
}
