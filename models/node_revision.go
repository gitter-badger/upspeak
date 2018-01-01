package models

import (
	"database/sql"
	"log"
	"time"
)

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

///////////////////////
// Add node revision //
///////////////////////
type AddNodeRevisionSchema struct {
	NodeID     int64
	RevisionID int64
	Subject    string
	Body       string
	UserID     int64
}

var addNodeRevisionQuery = `with rev as (
  insert into node_revisions
     (node_id, subject, body, committer_id)
  values (
    $1, -- NodeID
    $2, -- Subject
    $3, -- Body
    $4)  -- UserID
  returning id as rev_id, node_id
)
update nodes
  set revision_head = rev.rev_id
from rev where nodes.id = rev.node_id
returning rev.*;`

// GetNodeRevision gets a node revision
func AddNodeRevision(n *AddNodeRevisionSchema) (*AddNodeRevisionSchema, error) {
	err := db.QueryRow(getNodeRevisionQuery, &n.NodeID, &n.Subject, &n.Body, &n.UserID).Scan(&n.RevisionID, &n.NodeID)

	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil

}
