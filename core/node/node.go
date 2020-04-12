package node

import (
	"log"

	"github.com/upspeak/upspeak/core"
)

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

// Get returns a Node given a Node ID
func Get(nodeID *int64) (core.Node, error) {
	n := core.NewNode()
	var replyTo core.NullInt64
	var threadID core.NullInt64
	var updatedBy core.NullInt64
	err := core.DB.QueryRow(getNodeQuery, nodeID).Scan(
		&n.NodeID, &threadID, &replyTo,
		&n.Data.DataType, &n.Data.Subject, &n.Data.Body, &n.Data.RichData,
		&n.Meta.CreatedAt, &n.Meta.CreatedBy.UserID, &n.Meta.CreatedBy.Username,
		&n.Meta.UpdatedAt, &updatedBy,
	)

	if err != nil {
		log.Println(err)
		return n, err
	}
	if updatedBy.Valid {
		n.Meta.UpdatedBy = core.NodeAuthor{
			UserID: updatedBy,
		}
	}
	if threadID.Valid {
		n.Thread = &core.Thread{
			ThreadID: threadID,
		}
	}
	if replyTo.Valid {
		n.InReplyTo = &core.Node{
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

// Forks returns forks of a node
func Forks(nodeID int64) ([]*core.Thread, error) {
	rows, err := core.DB.Query(getNodeForksQuery, nodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	forks := make([]*core.Thread, 0)

	for rows.Next() {
		fork := new(core.Thread)
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

// Replies lists nodes that were replies made directly to a given node. This
// is not the same as listing child nodes of a thread.
func Replies(nodeID int64) ([]*core.Node, error) {
	rows, err := core.DB.Query(getRepliesQuery, nodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	replies := make([]*core.Node, 0)

	for rows.Next() {
		n := core.NewNode()
		var replyTo core.NullInt64
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
		if replyTo.Valid {
			n.InReplyTo = &core.Node{
				NodeID: replyTo,
			}
		}
		replies = append(replies, &n)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return replies, err
	}
	return replies, nil
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

// Fork forks a node into a thread and returns the new thread ID
func Fork(srcNodeID int64, targetTeamID int64, authorID int64, quotedData *core.NodeData) (int64, error) {
	var threadID int64
	err := core.DB.QueryRow(forkNodeQuery,
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
