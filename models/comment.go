package models

import (
	"log"
	"time"
)

////////////////////
// Create comment //
////////////////////
type CreateCommentSchema struct {
	NodeID      int64
	RevisionID  int64
	CreatedAt   time.Time
	ThreadID    int64
	InReplyToID int64
	DataType    string
	Subject     string
	Body        string
	Extra       JSONB
	UserID      int64
}

var createCommentQuery = `with ids as (
    -- Generate IDs before hand for nodes and node_revision
    select
        generate_id() as node_id,
        generate_id('node_revision_seq') as revision_id
),
n as (
    -- Insert node first
    insert into nodes
    (id, author_id, source_node_id, in_reply_to, data_type, revision_head, created_at)
    select
        node_id,
        user_id,
        $1, -- Thread ID
        $2, -- Node to which this node is a direct reply to. Can be null
        $3, -- default should be 'markdown',
        revision_id,
        now()
    from ids
)
-- Insert node revision next
insert into node_revisions(id, node_id, subject, body, extra, committer_id)
    select
        revision_id,
        node_id,
        $4, -- Subject
        $5, -- Body
        $6, -- Node revision extra depending on data_type
        $7 -- User ID
    from ids
returning node_id, id as rev_id, created_at; -- return the node id, new revision id and created_at; -- return the node id, new revision id and created_at`

// CreateComment creates a comment
func CreateComment(c *CreateCommentSchema) (*CreateCommentSchema, error) {
	err := db.QueryRow(createCommentQuery, &c.ThreadID, &c.InReplyToID, &c.DataType, &c.Subject, &c.Body, &c.Extra, &c.UserID).Scan(&c.NodeID, &c.RevisionID, &c.CreatedAt)
	if err != nil {
		log.Println(err)
		return c, err
	}
	return c, nil
}
