package models

import (
	"log"
	"time"
)

///////////////////////
// Type declarations //
///////////////////////

// NodeAuthor represents author of a specific node or its edits
type NodeAuthor struct {
	ID       *int64  `json:"id"`
	Username *string `json:"username,omitempty"`
}

// NodeData holds the content for a node
type NodeData struct {
	DataType *string `json:"data_type"`
	Subject  *string `json:"subject"`
	Body     *string `json:"body"`
	RichData *JSONB  `json:"rich_data"`
}

// Node represents a single node structure
type Node struct {
	ID        int64       `json:"id"`
	ThreadID  *int64      `json:"thread_id,omitempty"`
	Author    *NodeAuthor `json:"author"`
	Data      *NodeData   `json:"data"`
	CreatedAt *time.Time  `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at,omitempty"`
	UpdatedBy *int64      `json:"updated_by,omitempty"`
	InReplyTo *int64      `json:"in_reply_to,omitempty"`
}

// newNode returns an empty `Node` type which can be used to fill data
func newNode() *Node {
	n := new(Node)
	n.Author = new(NodeAuthor)
	n.Data = new(NodeData)
	return n
}

// Thread holds a list of nodes and some of its metadata
type Thread struct {
	ThreadID    *int64  `json:"thread_id,omitempty"`
	TeamID      *int64  `json:"team_id,omitempty"`
	SourceNode  *Node   `json:"source_node,omitempty"`
	ChildNodes  []*Node `json:"child_nodes,omitempty"`
	IsOpen      *bool   `json:"is_open,omitempty"`
	Permissions *JSONB  `json:"permissions,omitempty"`
	Attrs       *JSONB  `json:"attrs,omitempty"`
}

// NodeCreateRes is a standard response for creating different types of nodes
type NodeCreateRes struct {
	NodeID    int64      `json:"node_id"`
	ThreadID  *int64     `json:"thread_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

///////////////
// Functions //
///////////////

///////////////
// Get nodes //
///////////////

var getNodesQuery = `
select
	nodes.id as node_id, in_reply_to,
	data_type, subject, body, rich_data,
	nodes.created_at, author_id, users.username,
	updated_at, updated_by
from nodes
    join users on (nodes.author_id = users.id)
where nodes.source_node_id = $1
    or nodes.id = $1 -- this adds the source node detail without doing a union
order by nodes.id asc;
`

// GetNodes queries db for multiple nodes based on thread ID
func GetNodes(threadID *int64) ([]*Node, error) {
	rows, err := db.Query(getNodesQuery, threadID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	nodes := make([]*Node, 0)

	for rows.Next() {
		n := newNode()
		err := rows.Scan(
			&n.ID, &n.InReplyTo,
			&n.Data.DataType, &n.Data.Subject, &n.Data.Body, &n.Data.RichData,
			&n.CreatedAt, &n.Author.ID, &n.Author.Username,
			&n.UpdatedAt, &n.UpdatedBy,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		nodes = append(nodes, n)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nodes, err
	}

	return nodes, nil
}

//////////////
// Get node //
//////////////

var getNodeQuery = `
select
	nodes.id as node_id, source_node_id, in_reply_to,
	data_type, subject, body, rich_data,
	nodes.created_at, author_id, users.username,
	updated_at as last_edited_at, updated_by as last_edited_by
from nodes
    join users on (nodes.author_id = users.id)
where nodes.id = $1;
`

// GetNode returns a Node given a Node ID
func GetNode(nodeID *int64) (*Node, error) {
	n := newNode()
	err := db.QueryRow(getNodeQuery, nodeID).Scan(
		&n.ID, &n.ThreadID, &n.InReplyTo,
		&n.Data.DataType, &n.Data.Subject, &n.Data.Body, &n.Data.RichData,
		&n.CreatedAt, &n.Author.ID, &n.Author.Username,
		&n.UpdatedAt, &n.UpdatedBy,
	)
	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil
}

///////////////////////
// Get forks of node //
///////////////////////

var getNodeForksQuery = `select
    id, team_id, permissions, is_open, attrs
from threads
    where forked_from_node = $1 -- This is the current node ID
order by id asc;`

func GetNodeForks(nodeID int64) ([]*Thread, error) {
	rows, err := db.Query(getNodeForksQuery, nodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	forks := make([]*Thread, 0)

	for rows.Next() {
		fork := new(Thread)
		err := rows.Scan(&fork.ThreadID, &fork.TeamID, &fork.Permissions, &fork.IsOpen, &fork.Attrs)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		forks = append(forks, fork)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return forks, err
	}
	return forks, nil
}

///////////////////
// Create thread //
///////////////////

type CreateThreadSchema struct {
	TeamID     int64
	UserID     int64
	Data       NodeData
	NodeID     int64
	RevisionID int64
	CreatedAt  time.Time
}

var createThreadQuery = `
-- Insert node first
with n as (
    insert into nodes (author_id, data_type, subject, body, rich_data, created_at)
        values($1, $2, $3, $4, $5, now())
        returning id, created_at
)
-- insert thread
insert into threads(id, team_id)
	select id, $6 from n
	returning id as thread_id;
`

// CreateThread creates a node of type thread
func CreateThread(t *CreateThreadSchema) (int64, error) {
	var threadID int64
	err := db.QueryRow(createThreadQuery, &t.UserID, &t.Data.DataType, &t.Data.Subject, &t.Data.Body, &t.Data.RichData, &t.TeamID).Scan(&threadID)
	if err != nil {
		log.Println(err)
		return threadID, err
	}
	return threadID, nil
}

/////////////////
// Get replies //
/////////////////

var getRepliesQuery = `
select
	nodes.id,
	data_type, subject, body, rich_data,
	nodes.created_at, author_id, users.username,
	updated_at, updated_by
from nodes
    join users on (nodes.author_id = users.id)
where nodes.in_reply_to = $1
order by nodes.id asc;
`

// GetReplies lists nodes that were replies made directly to a given node. This
// is not the same as listing child nodes of a thread.
func GetReplies(nodeID int64) ([]*Node, error) {
	rows, err := db.Query(getRepliesQuery, nodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	replies := make([]*Node, 0)

	for rows.Next() {
		r := newNode()
		err := rows.Scan(
			&r.ID,
			&r.Data.DataType, &r.Data.Subject, &r.Data.Body, &r.Data.RichData,
			&r.CreatedAt, &r.Author.ID, &r.Author.Username,
			&r.UpdatedAt, &r.UpdatedBy,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		replies = append(replies, r)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return replies, err
	}
	return replies, nil
}

///////////////////////////
// Get forks in a thread //
///////////////////////////

type GetForksInAThreadSchema struct {
	NodeID      int64
	ThreadID    int64
	TeamID      int64
	Permissions JSONB
	IsOpen      bool
	Attrs       JSONB
}

var getForksInAThreadQuery = `select
    forked_from_node as node_id, -- The node which was forked
    id as forked_thread_id, -- The forked thread created
    team_id, permissions, is_open, attrs
from threads
    where forked/_from_node in (
        select id from nodes
            where source_node_id = $1 or id = $1 -- ThreadID
    )
order by forked_from_node, id;`

func GetForksInAThread(ThreadID int64) ([]*GetForksInAThreadSchema, error) {
	rows, err := db.Query(getForksInAThreadQuery, ThreadID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	forks := make([]*GetForksInAThreadSchema, 0)

	for rows.Next() {
		fork := new(GetForksInAThreadSchema)
		err := rows.Scan(&fork.NodeID, &fork.ThreadID, &fork.Permissions, &fork.IsOpen, &fork.Attrs)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		forks = append(forks, fork)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return forks, err
	}
	return forks, nil
}

///////////////
// Fork node //
///////////////

type ForkNodeSchema struct {
	NodeID     int64
	RevisionID int64
	TeamID     int64
	UserID     int64
	CreatedAt  time.Time
}

var forkNodeQuery = `with ids as (
    -- Generate IDs before hand for nodes and node_revision
    select generate_id() as node_id,
        generate_id('node_revision_seq') as revision_id,

        -- These should be passed from application
        $1 as team_id, -- TeamID
        $2 as user_id -- UserID
),
n as (
    -- Insert node first
    insert into nodes (id, author_id, data_type, revision_head, created_at)
        select node_id, user_id, 'markdown', revision_id, now() from ids
        returning created_at
),
rev as (
    -- Insert node revision next
    insert into node_revisions(id, node_id, subject, body, committer_id)
        -- :subject and :body should be set only if there is any new content added
        select revision_id, node_id, $3, $4, user_id from ids
), -- Subject and Body
thread as (
    -- insert thread
    insert into threads(id, team_id, forked_from_node)
        -- set :forked_from to node ID of original node that is being forked
        select node_id, team_id, $5 from ids -- ForkedFromID
)
-- return the result
select ids.*, n.* from ids, n;`

// ForkNode forks a node into a thread
func ForkNode(TeamID int64, UserID int64, Subject string, Body string, ForkedFrom int64) (*ForkNodeSchema, error) {
	thread := new(ForkNodeSchema)
	err := db.QueryRow(forkNodeQuery, TeamID, UserID, Subject, Body, ForkedFrom).Scan(&thread.NodeID, &thread.RevisionID, &thread.TeamID, &thread.UserID, &thread.CreatedAt)
	if err != nil {
		log.Println(err)
		return thread, err
	}
	return thread, nil
}
