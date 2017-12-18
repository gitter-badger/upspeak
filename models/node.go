package models

import (
	"database/sql"
	"log"
	"time"
)

type GetNodeSchema struct {
	NodeId       int64
	ThreadId     int64
	SourceNodeId int64
	AuthorId     int64
	Username     string
	CommitterId  int64
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

func GetNodes(n *GetNodeSchema) ([]*GetNodeSchema, error) {
	rows, err := db.Query(getNodesQuery, n.ThreadId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(rows)
	defer rows.Close()

	nodes := make([]*GetNodeSchema, 0)

	for rows.Next() {
		node := new(GetNodeSchema)
		err := rows.Scan(&n.NodeId, &n.AuthorId, &n.Username, &n.CommitterId, &n.DataType, &n.RevisionHead, &n.Subject, &n.Body, &n.Extra, &n.CreatedAt, &n.LastEditedAt)
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
	err := db.QueryRow(getNodeQuery, &n.NodeId).Scan(&n.NodeId, &n.SourceNodeId, &n.InReplyTo, &n.AuthorId, &n.Username, &n.CommitterId, &n.DataType, &n.RevisionHead, &n.Subject, &n.Body, &n.Extra, &n.CreatedAt, &n.LastEditedAt)
	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil
}
