package models

import (
	"database/sql"
	"log"
	"time"
)

type NodeRevision struct {
	NodeID    int64     `json:"node_id"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int64     `json:"updated_by"`
	Data      *NodeData `json:"data,omitempty"`
}

type GetNodeRevisionSchema struct {
	NodeID      int64
	RevisionID  int64
	CommitterID int64
	Username    string
	DataType    sql.NullString
	Subject     sql.NullString
	Body        sql.NullString
	Extra       JSONB
	CommittedAt time.Time
}

///////////////////////
// Get node revision //
///////////////////////
var getNodeRevisionQuery = `select
    node_revisions.id as revision_id, node_id,
    committer_id, users.username, -- User who created the revision
    data_type, subject, body, extra,
    node_revisions.created_at as committed_at -- When revision was created
from node_revisions
    join nodes on (node_revisions.node_id = nodes.id)
    join users on (node_revisions.committer_id = users.id)
where node_revisions.node_id = $1
order by node_revisions.created_at desc; -- latest first
`

// GetNodeRevision gets a node revision
func GetNodeRevision(n *GetNodeRevisionSchema) (*GetNodeRevisionSchema, error) {
	err := db.QueryRow(getNodeRevisionQuery, &n.NodeID).Scan(&n.RevisionID, &n.CommitterID, &n.Username, &n.DataType, &n.Subject, &n.Body, &n.Extra, &n.CommittedAt)

	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil

}

////////////////////////
// Get node revisions //
////////////////////////
var getNodeRevisionsQuery = `select
    node_revisions.id as revision_id, node_id,
    committer_id, users.username, -- User who created the revision
    data_type, subject, body, extra,
    node_revisions.created_at as committed_at -- When revision was created
from node_revisions
    join nodes on (node_revisions.node_id = nodes.id)
    join users on (node_revisions.committer_id = users.id)
where node_revisions.node_id = $1
order by node_revisions.created_at desc; -- latest first`

// GetNodeRevisions gets node revisions
func GetNodeRevisions(NodeID int64) ([]*GetNodeRevisionSchema, error) {
	rows, err := db.Query(getNodeRevisionsQuery, NodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	nodeRevisions := make([]*GetNodeRevisionSchema, 0)

	for rows.Next() {
		nodeRevision := new(GetNodeRevisionSchema)
		err := rows.Scan(&nodeRevision.RevisionID, &nodeRevision.CommitterID, &nodeRevision.Username, &nodeRevision.DataType, &nodeRevision.Subject, &nodeRevision.Body, &nodeRevision.Extra, &nodeRevision.CommittedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		nodeRevisions = append(nodeRevisions, nodeRevision)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nodeRevisions, err
	}
	return nodeRevisions, nil
}

////////////////////////
// Edit a node's data //
////////////////////////

type NodeAddRevisionSchema struct {
	NodeID    int64
	UpdatedBy int64
	Subject   string
	Body      string
	RichData  JSONB
}

var nodeAddRevisionQuery = `
update nodes set
	subject = $1,
	body = $2,
	rich_data = $3,
	updated_by = $4,
	updated_at = now()
where id = $5
returning id, updated_at;
`

// NodeAddRevision updates the contents of a node and internally creates a new revision
func NodeAddRevision(n *NodeAddRevisionSchema) (*NodeRevision, error) {
	r := new(NodeRevision)
	err := db.QueryRow(
		nodeAddRevisionQuery,
		&n.Subject, &n.Body, &n.RichData, &n.UpdatedBy, &n.NodeID,
	).Scan(&r.NodeID, &r.UpdatedAt)

	if err != nil {
		log.Println(err)
		return r, err
	}
	r.UpdatedBy = n.UpdatedBy
	return r, nil

}
