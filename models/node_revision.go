package models

import (
	"database/sql"
	"log"
	"time"
)

type NodeRevision struct {
	CreatedAt time.Time  `json:"created_at"`
	Committer NodeAuthor `json:"committer"`
	Data      NodeData   `json:"data,omitempty"`
}

type GetNodeRevisionSchema struct {
	NodeID      int64
	CommittedAt time.Time
}

///////////////////////
// Get node revision //
///////////////////////
var getNodeRevisionQuery = `
select
    r.committer_id, users.username,
    nodes.data_type, r.subject, r.body, r.rich_data,
    r.created_at
from audit.node_revisions r
	join public.nodes on (r.node_id = public.nodes.id)
    join public.users on (r.committer_id = public.users.id)
where r.node_id = $1 and r.created_at = $2;
`

// GetNodeRevision gets a node revision
func GetNodeRevision(n *GetNodeRevisionSchema) (*NodeRevision, error) {
	r := new(NodeRevision)
	err := db.QueryRow(getNodeRevisionQuery, &n.NodeID, &n.CommittedAt).Scan(
		&r.Committer.ID, &r.Committer.Username,
		&r.Data.DataType, &r.Data.Subject, &r.Data.Body, &r.Data.RichData,
		&r.CreatedAt,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		return nil, err
	}
	return r, nil
}

////////////////////////
// Get node revisions //
////////////////////////
var getNodeRevisionsQuery = `
select
    r.committer_id, users.username,
    nodes.data_type, r.subject, r.body, r.rich_data,
    r.created_at
from audit.node_revisions r
	join public.nodes on (r.node_id = public.nodes.id)
    join public.users on (r.committer_id = public.users.id)
where r.node_id = $1
order by r.created_at desc; -- latest first
`

// GetNodeRevisions gets node revisions
func GetNodeRevisions(NodeID int64) ([]*NodeRevision, error) {
	rows, err := db.Query(getNodeRevisionsQuery, NodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	nodeRevisions := make([]*NodeRevision, 0)

	for rows.Next() {
		r := new(NodeRevision)
		err := rows.Scan(
			&r.Committer.ID, &r.Committer.Username,
			&r.Data.DataType, &r.Data.Subject, &r.Data.Body, &r.Data.RichData,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		nodeRevisions = append(nodeRevisions, r)
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
returning updated_at;
`

// NodeAddRevision updates the contents of a node and internally creates a new revision
func NodeAddRevision(n *NodeAddRevisionSchema) (*NodeRevision, error) {
	r := new(NodeRevision)
	err := db.QueryRow(
		nodeAddRevisionQuery,
		&n.Subject, &n.Body, &n.RichData, &n.UpdatedBy, &n.NodeID,
	).Scan(&r.CreatedAt)

	if err != nil {
		log.Println(err)
		return r, err
	}
	r.Committer.ID = n.UpdatedBy
	return r, nil

}
