package thread

import (
	"log"

	"github.com/upspeak/upspeak/core"
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
returning id;
`

// AddComment creates a comment and returns the new node ID
func AddComment(data core.NodeData, threadID int64, authorID int64, inReplyToID int64) (int64, error) {
	var nodeID int64
	err := core.DB.QueryRow(
		createCommentQuery,
		authorID, threadID, inReplyToID,
		data.DataType, data.Subject, data.Body, data.RichData,
	).Scan(&nodeID)
	if err != nil {
		log.Println(err)
		return nodeID, err
	}
	return nodeID, nil
}
