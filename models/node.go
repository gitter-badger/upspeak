package models

import (
	"database/sql"
	"log"
	"time"
)

type GetNodeSchema struct {
	NodeID       int64
	ThreadID     int64
	SourceNodeID sql.NullInt64
	AuthorID     int64
	Username     string
	CommitterID  int64
	InReplyTo    sql.NullString
	DataType     string
	RevisionHead int64
	Subject      sql.NullString
	Body         string
	CreatedAt    time.Time
	Extra        sql.NullString
	LastEditedAt time.Time
}

///////////////
// Get nodes //
///////////////
var getNodesQuery = `select
    node_id, author_id, users.username,
    committer_id, data_type, revision_head,
    subject, body, extra,
    nodes.created_at, rev.created_at as last_edited_at
    -- add further columns as needed
from nodes
    join node_revisions rev on (nodes.revision_head = rev.id)
    join users on (nodes.author_id = users.id)
where nodes.source_node_id = $1
    or nodes.id = $1 -- this adds the source node detail without doing a union
order by nodes.created_at asc, nodes.id asc;`

func GetNodes(ThreadID int64) ([]*GetNodeSchema, error) {
	rows, err := db.Query(getNodesQuery, ThreadID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(rows)
	defer rows.Close()

	nodes := make([]*GetNodeSchema, 0)

	for rows.Next() {
		node := new(GetNodeSchema)
		err := rows.Scan(&node.NodeID, &node.AuthorID, &node.Username, &node.CommitterID, &node.DataType, &node.RevisionHead, &node.Subject, &node.Body, &node.Extra, &node.CreatedAt, &node.LastEditedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		nodes = append(nodes, node)
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
var getNodeQuery = `select
    node_id, source_node_id, in_reply_to,
    author_id, users.username,
    committer_id, data_type, revision_head,
    subject, body, extra,
    nodes.created_at, rev.created_at as last_edited_at
    -- add further columns as needed
from nodes
    join node_revisions rev on (nodes.revision_head = rev.id)
    join users on (nodes.author_id = users.id)
where nodes.id = $1;`

func GetNode(n *GetNodeSchema) (*GetNodeSchema, error) {
	err := db.QueryRow(getNodeQuery, &n.NodeID).Scan(&n.NodeID, &n.SourceNodeID, &n.InReplyTo, &n.AuthorID, &n.Username, &n.CommitterID, &n.DataType, &n.RevisionHead, &n.Subject, &n.Body, &n.Extra, &n.CreatedAt, &n.LastEditedAt)
	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil
}

///////////////////////
// Get forks of node //
///////////////////////
type GetForksOfNodeSchema struct {
	NodeID      int64
	ThreadID    int64
	TeamID      int64
	Permissions JSONB
	IsOpen      bool
	Attrs       JSONB
}

var getForksOfNodeQuery = `select
    id, team_id, permissions, is_open, attrs
from threads
    where forked_from_node = $1 -- This is the current node ID
order by id asc;`

func GetForksOfNode(NodeID int64) ([]*GetForksOfNodeSchema, error) {
	rows, err := db.Query(getForksOfNodeQuery, NodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	forks := make([]*GetForksOfNodeSchema, 0)

	for rows.Next() {
		fork := new(GetForksOfNodeSchema)
		err := rows.Scan(&fork.NodeID, &fork.TeamID, &fork.Permissions, &fork.IsOpen, &fork.Attrs)
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
	Subject    string
	Body       string
	NodeID     int64
	RevisionID int64
	CreatedAt  time.Time
}

var createThreadQuery = `with ids as (
-- Generate IDs before hand for nodes and node_revision

 select generate_id() as node_id,
        generate_id('node_revision_seq') as revision_id,

        -- These should be passed from application
        $1 as team_id,
        $2 as user_id
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
        select revision_id, node_id, $3, $4, user_id from ids
),
thread as (
    -- insert thread
    insert into threads(id, team_id)
        select node_id, team_id from ids
)
-- return the result
select ids.*, n.* from ids, n;`

// CreateThread creates a node of type thread
func CreateThread(t *CreateThreadSchema) (*CreateThreadSchema, error) {
	err := db.QueryRow(createThreadQuery, &t.TeamID, &t.UserID, &t.Subject, &t.Body).Scan(&t.NodeID, &t.RevisionID, &t.TeamID, &t.UserID, &t.CreatedAt)
	if err != nil {
		log.Println(err)
		return t, err
	}
	return t, nil
}

/////////////////
// Get replies //
/////////////////
type GetRepliesSchema struct {
	NodeID       int64
	AuthorID     int64
	Username     string
	CommitterID  string
	DataType     sql.NullString
	RevisionHead int64
	Subject      sql.NullString
	Body         sql.NullString
	Extra        JSONB
	CreatedAt    time.Time
	LastEditedAt time.Time
}

var getRepliesQuery = `with n as (
    select id, source_node_id from nodes where id = $1 -- ThreadID
),
ns as (
    select
        n.id as reply_target,
        nodes.id, n.source_node_id, in_reply_to, author_id, data_type, revision_head, created_at
    from nodes, n where nodes.source_node_id = (
        -- Iterate to figure out the source node of given node.
        -- If given node is a source node itself, we use its id
        -- Else, we use the source_node_id for given node.
        case when n.source_node_id is null then
            n.id
        else
            n.source_node_id
        end
    )
)
select
    node_id, author_id, users.username,
    committer_id, data_type, revision_head,
    subject, extra,
    ns.created_at, rev.created_at as last_edited_at
from ns
    join node_revisions rev on (ns.revision_head = rev.id)
    join users on (ns.author_id = users.id)
where ns.in_reply_to = ns.reply_target
order by ns.id asc;`

func GetReplies(ThreadID int64) ([]*GetRepliesSchema, error) {
	rows, err := db.Query(getRepliesQuery, ThreadID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	replies := make([]*GetRepliesSchema, 0)

	for rows.Next() {
		reply := new(GetRepliesSchema)
		err := rows.Scan(&reply.NodeID, &reply.AuthorID, &reply.Username, &reply.CommitterID, &reply.DataType, &reply.RevisionHead, &reply.Subject, &reply.Body, &reply.Extra, &reply.CreatedAt, &reply.LastEditedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		replies = append(replies, reply)
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
