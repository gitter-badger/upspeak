package rpc

import (
	"database/sql"
	"log"
	"net/http"
	"time"

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
	TeamID   *int64        `json:"team_id"`
	AuthorID *int64        `json:"author_id"`
	DataType *string       `json:"data_type"`
	Subject  *string       `json:"subject"`
	Body     *string       `json:"body"`
	RichData *models.JSONB `json:"rich_data"`
}

type CreateThreadReply struct {
	ThreadID int64 `json:"thread_id"`
}

func (n *NodeService) CreateThread(r *http.Request, args *CreateThreadArgs, reply *CreateThreadReply) error {
	data := &models.NodeData{
		DataType: args.DataType,
		Subject:  args.Subject,
		Body:     args.Body,
		RichData: args.RichData,
	}
	threadID, err := models.CreateThread(data, args.TeamID, args.AuthorID)
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
	NodeID   *int64        `json:"node_id"`
	AuthorID *int64        `json:"author_id"`
	Subject  *string       `json:"subject"`
	Body     *string       `json:"body"`
	RichData *models.JSONB `json:"rich_data"`
}

type NodeEditReply struct {
	NodeID   *int64               `json:"node_id"`
	Revision *models.NodeRevision `json:"revision"`
}

// Edit creates a node revision
func (n *NodeService) Edit(r *http.Request, args *NodeEditArgs, reply *NodeEditReply) error {
	data := &models.NodeData{
		Subject:  args.Subject,
		Body:     args.Body,
		RichData: args.RichData,
	}
	rev, err := models.NodeAddRevision(args.NodeID, args.AuthorID, data)
	if err != nil {
		log.Println(err)
		return err
	}

	reply.NodeID = args.NodeID
	reply.Revision = rev
	return nil
}

/////////////////
// Add comment //
/////////////////

type NodeCreateCommentArgs struct {
	ThreadID    *int64        `json:"thread_id"`
	InReplyToID *int64        `json:"in_reply_to_id"`
	DataType    *string       `json:"data_type"`
	Subject     *string       `json:"subject"`
	Body        *string       `json:"body"`
	RichData    *models.JSONB `json:"rich_data"`
	AuthorID    *int64        `json:"author_id"`
}

type NodeCreateCommentReply struct {
	NodeID int64 `json:"node_id"`
}

// CreateComment adds a comment
func (n *NodeService) CreateComment(r *http.Request, args *NodeCreateCommentArgs, reply *NodeCreateCommentReply) error {
	data := &models.NodeData{
		DataType: args.DataType,
		Subject:  args.Subject,
		Body:     args.Body,
		RichData: args.RichData,
	}
	nid, err := models.CreateComment(data, args.ThreadID, args.AuthorID, args.InReplyToID)
	if err != nil {
		log.Println(err)
		return err
	}

	reply.NodeID = nid
	return nil
}

///////////////////////
// Get node revision //
///////////////////////

type NodeGetRevisionArgs struct {
	NodeID    *int64     `json:"node_id"`
	CreatedAt *time.Time `json:"created_at"`
}

type NodeGetRevisionReply struct {
	NodeID   *int64               `json:"node_id"`
	Revision *models.NodeRevision `json:"revision"`
}

// GetRevision gets a node revision
func (n *NodeService) GetRevision(r *http.Request, args *NodeGetRevisionArgs, reply *NodeGetRevisionReply) error {
	rev, err := models.GetNodeRevision(args.NodeID, args.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
	}

	// Generate response
	reply.NodeID = args.NodeID
	reply.Revision = rev
	return nil
}

////////////////////////
// Get node revisions //
////////////////////////

type NodeGetRevisionsArgs struct {
	NodeID *int64 `json:"node_id"`
}

type NodeGetRevisionsReply struct {
	NodeID        *int64                 `json:"node_id"`
	RevisionCount int                    `json:"revision_count"`
	Revisions     []*models.NodeRevision `json:"revisions"`
}

// GetRevisions gets revisions of a node
func (n *NodeService) GetRevisions(r *http.Request, args *NodeGetRevisionsArgs, reply *NodeGetRevisionsReply) error {
	revisions, err := models.GetNodeRevisions(args.NodeID)
	if err != nil {
		log.Println(err)
		return err
	}

	// Generate response
	reply.NodeID = args.NodeID
	reply.RevisionCount = len(revisions)
	reply.Revisions = revisions
	return nil

}

//////////////////////
// Get node replies //
//////////////////////

type NodeGetRepliesArgs struct {
	NodeID int64 `json:"node_id"`
}

type NodeGetRepliesReply struct {
	NodeID     *int64         `json:"node_id"`
	ReplyCount int            `json:"reply_count"`
	Replies    []*models.Node `json:"replies"`
}

// GetReplies gets replies of a node
func (n *NodeService) GetReplies(r *http.Request, args *NodeGetRepliesArgs, reply *NodeGetRepliesReply) error {
	nodeReplies, err := models.GetReplies(args.NodeID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.NodeID = &args.NodeID
	reply.ReplyCount = len(nodeReplies)
	reply.Replies = nodeReplies

	return nil
}

///////////////////////////
// Get forks in a thread //
///////////////////////////

type NodeGetForksInAThreadArgs struct {
	ThreadID int64 `json:"thread_id"`
}

type NodeGetForksInAThreadReply struct {
	ThreadID  int64            `json:"thread_id"`
	ForkCount int              `json:"fork_count"`
	Forks     []*models.Thread `json:"forks"`
}

// GetForksInThread returns list of forks of all nodes in the thread
func (n *NodeService) GetForksInAThread(r *http.Request, args *NodeGetForksInAThreadArgs, reply *NodeGetForksInAThreadReply) error {
	forks, err := models.GetForksInAThread(args.ThreadID)
	if err != nil {
		log.Println(err)
	}

	// Generate response
	reply.ThreadID = args.ThreadID
	reply.ForkCount = len(forks)
	reply.Forks = forks
	return nil
}

///////////////
// Fork node //
///////////////

type NodeForkNodeArgs struct {
	NodeID   *int64        `json:"node_id"`
	TeamID   *int64        `json:"team_id"`
	AuthorID *int64        `json:"author_id"`
	Subject  *string       `json:"subject"`
	Body     *string       `json:"body"`
	RichData *models.JSONB `json:"rich_data"`
}

type NodeForkNodeReply struct {
	ForkedThreadID *int64 `json:"forked_thread_id"`
	ForkedFrom     *int64 `json:"forked_from"`
}

// Fork creates a new thread from current node
func (n *NodeService) Fork(r *http.Request, args *NodeForkNodeArgs, reply *NodeForkNodeReply) error {
	quotedData := &models.NodeData{
		Subject:  args.Subject,
		Body:     args.Body,
		RichData: args.RichData,
	}
	threadID, err := models.ForkNode(args.NodeID, args.TeamID, args.AuthorID, quotedData)
	if err != nil {
		log.Println(err)
		return err
	}

	// Generate response
	reply.ForkedThreadID = &threadID
	reply.ForkedFrom = args.NodeID
	return nil
}

////////////////
// Node utils //
////////////////

// resolvePermissions returns node's permission matrix
func (n *NodeService) resolvePermissions(nodeID int64) {
	// TODO: Define permission matrix data structure
}
