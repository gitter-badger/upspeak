package models

import (
	"log"
	"time"
)

type GetNodeRevisionSchema struct {
	NodeId      int64
	RevisionId  int64
	CommitterId int64
	Username    string
	DataType    string
	Subject     string
	Body        string
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
	err := db.QueryRow(getNodeRevisionQuery, n.NodeId).Scan(n.RevisionId, n.CommitterId, n.Username, n.DataType, n.Subject, n.Body, n.Extra, n.CommittedAt)

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
func GetNodeRevisions(n *GetNodeRevisionSchema) (*GetNodeRevisionSchema, error) {
	rows, err := db.Query(getNodeRevisionsQuery, n.NodeId)
	if err != nil {
		log.Println(err)
		return n, err
	}
	defer rows.Close()

	nodeRevisions := make([]*GetNodeRevisionSchema, 0)

	for rows.Next() {
		nodeRevision := new(GetNodeRevisionSchema)
		err := rows.Scan(nodeRevision.RevisionId, nodeRevision.CommitterId, nodeRevision.Username, nodeRevision.DataType, nodeRevision.Subject, nodeRevision.Body, nodeRevision.Extra, nodeRevision.CommittedAt)
		if err != nil {
			log.Println(err)
			return n, err
		}
		nodeRevisions = append(nodeRevisions, nodeRevision)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return n, err
	}
	return n, nil
}
