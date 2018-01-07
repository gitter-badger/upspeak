package models

import (
	"log"
)

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
		var replyTo *int64
		var updatedBy *int64
		err := rows.Scan(
			&n.NodeID, &replyTo,
			&n.Data.DataType, &n.Data.Subject, &n.Data.Body, &n.Data.RichData,
			&n.Meta.CreatedAt, &n.Meta.CreatedBy.UserID, &n.Meta.CreatedBy.Username,
			&n.Meta.UpdatedAt, &updatedBy,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		if updatedBy != nil {
			n.Meta.UpdatedBy = &NodeAuthor{
				UserID: updatedBy,
			}
		}
		if replyTo != nil {
			n.InReplyTo = &Node{
				NodeID: replyTo,
			}
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
	var replyTo *int64
	var threadID *int64
	var updatedBy *int64
	err := db.QueryRow(getNodeQuery, nodeID).Scan(
		&n.NodeID, &threadID, &replyTo,
		&n.Data.DataType, &n.Data.Subject, &n.Data.Body, &n.Data.RichData,
		&n.Meta.CreatedAt, &n.Meta.CreatedBy.UserID, &n.Meta.CreatedBy.Username,
		&n.Meta.UpdatedAt, &updatedBy,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	if updatedBy != nil {
		n.Meta.UpdatedBy = &NodeAuthor{
			UserID: updatedBy,
		}
	}
	if threadID != nil {
		n.Thread = &Thread{
			ThreadID: threadID,
		}
	}
	if replyTo != nil {
		n.InReplyTo = &Node{
			NodeID: replyTo,
		}
	}
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

// GetNodeForks returns forks of a node
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
func CreateThread(data *NodeData, teamID *int64, authorID *int64) (int64, error) {
	var threadID int64
	err := db.QueryRow(createThreadQuery,
		&authorID, &data.DataType, &data.Subject, &data.Body, &data.RichData, &teamID).Scan(&threadID)
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
		n := newNode()
		var replyTo *int64
		err := rows.Scan(
			&n.NodeID, &replyTo,
			&n.Data.DataType, &n.Data.Subject, &n.Data.Body, &n.Data.RichData,
			&n.Meta.CreatedAt, &n.Meta.CreatedBy.UserID, &n.Meta.CreatedBy.Username,
			&n.Meta.UpdatedAt, &n.Meta.UpdatedBy.UserID,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		if replyTo != nil {
			n.InReplyTo = &Node{
				NodeID: replyTo,
			}
		}
		replies = append(replies, n)
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

var getForksInAThreadQuery = `select
    forked_from_node as node_id, -- The node which was forked
    id as forked_thread_id, -- The forked thread created
    team_id, permissions, is_open, attrs
from threads
    where forked_from_node in (
        select id from nodes
            where source_node_id = $1 or id = $1 -- ThreadID
    )
order by forked_from_node, id;`

// GetForksInAThread returns list of all forks of all nodes in a thread
func GetForksInAThread(threadID int64) ([]*Thread, error) {
	rows, err := db.Query(getForksInAThreadQuery, threadID)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	forks := make([]*Thread, 0)

	for rows.Next() {
		fork := new(Thread)
		err := rows.Scan(&fork.ForkedFrom, &fork.ThreadID, &fork.TeamID, &fork.Permissions, &fork.IsOpen, &fork.Attrs)
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

var forkNodeQuery = `
-- Insert node first
with src as (
	select data_type from nodes where id = $1
),
n as (
    insert into nodes (author_id, data_type, subject, body, rich_data, created_at)
        select $2, data_type, $3, $4, $5, now() from src
        returning id as node_id
)
-- insert thread
insert into threads(id, team_id, forked_from_node)
	select node_id, $6, $1 from n
	returning id as thread_id;
`

// ForkNode forks a node into a thread and returns the new thread ID
func ForkNode(srcNodeID *int64, targetTeamID *int64, authorID *int64, quotedData *NodeData) (int64, error) {
	var threadID int64
	err := db.QueryRow(forkNodeQuery,
		srcNodeID, authorID,
		quotedData.Subject, quotedData.Body, quotedData.RichData,
		targetTeamID,
	).Scan(&threadID)
	if err != nil {
		log.Println(err)
		return threadID, err
	}
	return threadID, nil
}
