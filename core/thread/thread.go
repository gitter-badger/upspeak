package thread

import (
	"log"

	"github.com/upspeak/upspeak/core"
)

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

// New creates a new thread in DB with a new source node.
func New(data *core.NodeData, teamID *int64, authorID *int64) (int64, error) {
	var threadID int64
	err := core.DB.QueryRow(createThreadQuery,
		&authorID, &data.DataType, &data.Subject, &data.Body, &data.RichData, &teamID).Scan(&threadID)
	if err != nil {
		log.Println(err)
		return threadID, err
	}
	return threadID, nil
}

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

// Nodes queries db for multiple nodes based on thread ID.
// Use this to get all nodes of a thread
func Nodes(threadID int64) ([]*core.Node, error) {
	rows, err := core.DB.Query(getNodesQuery, threadID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	nodes := make([]*core.Node, 0)

	for rows.Next() {
		n := core.NewNode()
		var replyTo core.NullInt64
		var updatedBy core.NullInt64
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

		// If updatedBy is `not null` in db
		if updatedBy.Valid {
			n.Meta.UpdatedBy = core.NodeAuthor{
				UserID: updatedBy,
			}
		}

		if replyTo.Valid {
			n.InReplyTo = &core.Node{
				NodeID: replyTo,
			}
		}
		nodes = append(nodes, &n)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nodes, err
	}

	return nodes, nil
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

// Forks returns list of all forks of all nodes in a thread
func Forks(threadID int64) ([]*core.Thread, error) {
	rows, err := core.DB.Query(getForksInAThreadQuery, threadID)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	forks := make([]*core.Thread, 0)

	for rows.Next() {
		fork := new(core.Thread)
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
