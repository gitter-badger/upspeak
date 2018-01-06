package models

import (
	"log"
)

////////////////////
// Create comment //
////////////////////

var createCommentQuery = `
insert into nodes
    (
        author_id, source_node_id, in_reply_to,
        data_type, subject, body, rich_data,
        created_at
    )
values ($1, $2, $3, $4, $5, $6, $7, now())
returning id, created_at;
`

// CreateComment creates a comment
func CreateComment(c *Node) (*NodeCreateRes, error) {
	res := new(NodeCreateRes)
	err := db.QueryRow(
		createCommentQuery,
		c.Author.ID, c.ThreadID, c.InReplyTo, c.Data.DataType, c.Data.Subject, c.Data.Body, c.Data.RichData,
	).Scan(&res.NodeID, &res.CreatedAt)
	if err != nil {
		log.Println(err)
		return res, err
	}
	res.ThreadID = c.ThreadID
	return res, nil
}
