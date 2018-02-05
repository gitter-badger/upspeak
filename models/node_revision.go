package models

import (
	"database/sql"
	"log"
	"time"
)

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

// GetNodeRevision gets a node revision using node ID and timestamp of revision
func GetNodeRevision(nodeID int64, createdAt time.Time) (NodeRevision, error) {
	r := newNodeRevision()
	err := db.QueryRow(getNodeRevisionQuery, nodeID, createdAt).Scan(
		&r.Author.UserID, &r.Author.Username,
		&r.Data.DataType, &r.Data.Subject, &r.Data.Body, &r.Data.RichData,
		&r.CreatedAt,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		return r, err
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
func GetNodeRevisions(nodeID int64) ([]NodeRevision, error) {
	rows, err := db.Query(getNodeRevisionsQuery, nodeID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	nodeRevisions := make([]NodeRevision, 0)

	for rows.Next() {
		r := newNodeRevision()
		err := rows.Scan(
			&r.Author.UserID, &r.Author.Username,
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
func NodeAddRevision(nodeID NullInt64, authorID NullInt64, data *NodeData) (NodeRevision, error) {
	r := newNodeRevision()
	err := db.QueryRow(
		nodeAddRevisionQuery,
		&data.Subject, &data.Body, &data.RichData, &authorID, &nodeID,
	).Scan(&r.CreatedAt)

	if err != nil {
		log.Println(err)
		return r, err
	}
	r.Author = NodeAuthor{
		UserID: authorID,
	}
	return r, nil
}
