package models

import (
	"log"
	"time"
)

///////////////////
// Create thread //
///////////////////
type CreateThreadSchema struct {
	TeamId     int64
	UserId     int64
	Subject    string
	Body       string
	NodeId     int64
	RevisionId int64
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
	err := db.QueryRow(createThreadQuery, t.TeamId, t.UserId, t.Subject, t.Body).Scan(t.NodeId, t.RevisionId, t.TeamId, t.UserId, t.CreatedAt)
	if err != nil {
		log.Println(err)
		return t, err
	}
	return t, nil
}
