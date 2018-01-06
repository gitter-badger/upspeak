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

// NodeGetArgs defines the input parameters for node.Get
type NodeGetArgs struct {
	NodeID int64 `json:"node_id"`
}

// NodeGetReply defines the result data type for node.Get
type NodeGetReply struct {
	Node *models.Node `json:"node"`
}

// Get fetches the contents of a single node given a Node ID
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

// NodeGetForksArgs defines the input parameters for node.GetForks
type NodeGetForksArgs struct {
	NodeID int64 `json:"node_id"`
}

// NodeGetForksReply defines the result data type for node.GetForks
type NodeGetForksReply struct {
	NodeID    int64            `json:"node_id"`
	ForkCount int              `json:"fork_count"`
	Forks     []*models.Thread `json:"forks"`
}

// GetForks returns threads that were forked from a given node
func (n *NodeService) GetForks(r *http.Request, args *NodeGetForksArgs, reply *NodeGetForksReply) error {
	ts, err := models.GetNodeForks(args.NodeID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.NodeID = args.NodeID
	reply.ForkCount = len(ts)
	reply.Forks = ts

	return nil
}

//////////////////////////
// List nodes of thread //
//////////////////////////

// GetNodesArgs defines input parameters for node.GetNodes
type GetNodesArgs struct {
	ThreadID int64 `json:"thread_id"`
}

// GetNodesReply defines result data type for node.GetNodes
type GetNodesReply struct {
	ThreadID int64          `json:"thread_id"`
	Count    int            `json:"count"`
	Nodes    []*models.Node `json:"nodes"`
}

// GetNodes lists nodes for a given thread
func (n *NodeService) GetNodes(r *http.Request, args *GetNodesArgs, reply *GetNodesReply) error {
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
	NodeID   int64        `json:"node_id"`
	UserID   int64        `json:"user_id"`
	Subject  string       `json:"subject"`
	Body     string       `json:"body"`
	RichData models.JSONB `json:"rich_data"`
}

type NodeEditReply struct {
	Revision *models.NodeRevision `json:"revision"`
}

// Edit creates a node revision
func (n *NodeService) Edit(r *http.Request, args *NodeEditArgs, reply *NodeEditReply) error {
	rev, err := models.NodeAddRevision(&models.NodeAddRevisionSchema{
		NodeID:    args.NodeID,
		UpdatedBy: args.UserID,
		Subject:   args.Subject,
		Body:      args.Body,
		RichData:  args.RichData,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	reply.Revision = rev
	return nil
}

/////////////////
// Add comment //
/////////////////

type NodeCreateCommentArgs struct {
	ThreadID    models.NullInt64  `json:"thread_id"`
	InReplyToID models.NullInt64  `json:"in_reply_to_id"`
	DataType    string            `json:"data_type"`
	Subject     models.NullString `json:"subject"`
	Body        models.NullString `json:"body"`
	RichData    models.JSONB      `json:"rich_data"`
	UserID      int64             `json:"user_id"`
}

type NodeCreateCommentReply struct {
	Comment *models.NodeCreateRes `json:"comment"`
}

// CreateComment adds a comment
func (n *NodeService) CreateComment(r *http.Request, args *NodeCreateCommentArgs, reply *NodeCreateCommentReply) error {
	c, err := models.CreateComment(&models.Node{
		ThreadID:  &args.ThreadID,
		InReplyTo: &args.InReplyToID,
		Data: models.NodeData{
			DataType: args.DataType,
			Subject:  args.Subject,
			Body:     args.Body,
			RichData: args.RichData,
		},
		Author: models.NodeAuthor{
			ID: args.UserID,
		},
	})
	if err != nil {
		log.Println(err)
		return err
	}

	reply.Comment = c
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
