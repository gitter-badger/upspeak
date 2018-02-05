// Reusable types used and exported by models.
//
// These types are created with these goals:
//
// - Use composite types in models - complex types are composed of simple types.
//   See `Node` below for example.
// - Functions in the models package:
//    - return these types or any primitive type
//    - accept arguments of these types or use primitive types
//    - do not modify/mutate their arguments; return new copies or instances of a relevant type
//    - can be combined by higher order functions to produce extended outputs

package models

// NodeAuthor represents author of a specific node or its edits
type NodeAuthor struct {
	// ID of the user
	UserID NullInt64 `json:"user_id"`
	// Username of the user
	Username NullString `json:"username,omitempty"`
}

// NodeData holds the content for a node
type NodeData struct {
	// Data type of the node
	DataType NullString `json:"data_type"`
	// Subject field for the node
	Subject NullString `json:"subject"`
	// Body field for the node
	Body NullString `json:"body"`
	// RichData field for the node
	RichData JSONB `json:"rich_data"`
}

// NodeMeta contains information about node
type NodeMeta struct {
	// Time when node was created
	CreatedAt NullTime `json:"created_at"`
	// User who created the node
	CreatedBy NodeAuthor `json:"author"`
	// Time when node was last updated
	UpdatedAt NullTime `json:"updated_at,omitempty"`
	// User who last updated the node
	UpdatedBy NodeAuthor `json:"updated_by,omitempty"`
}

func newNodeMeta() NodeMeta {
	var n NodeMeta
	n.CreatedBy = NodeAuthor{}
	return n
}

// NodeRevision holds information about a specific node revision point
type NodeRevision struct {
	// Time when the revision was created
	CreatedAt NullTime `json:"created_at"`
	// User who created the revision
	Author NodeAuthor `json:"author"`
	// Data at this revision
	Data NodeData `json:"data,omitempty"`
}

func newNodeRevision() NodeRevision {
	var r NodeRevision
	r.Author = NodeAuthor{}
	r.Data = NodeData{}
	return r
}

// Node represents a single node structure
type Node struct {
	// Node ID
	NodeID NullInt64 `json:"node_id"`
	// Metadata information for the node
	Meta NodeMeta `json:"meta,omitempty"`
	// Details of the thread this node belongs to
	Thread *Thread `json:"thread_id,omitempty"`
	// Node's data
	Data NodeData `json:"data,omitempty"`
	// Node to which this node is a reply
	InReplyTo *Node `json:"in_reply_to,omitempty"`
	// Nodes which are replies to this node
	Replies []*Node `json:"replies,omitempty"`
	// Threads that have been forked from this node
	Forks []*Thread `json:"forks,omitempty"`
	// Revisions of the node
	Revisions []*NodeRevision `json:"revisions,omitempty"`
}

// newNode returns an empty `Node` type which can be used to fill data
func newNode() Node {
	var n Node
	n.Meta = newNodeMeta()
	n.Data = NodeData{}
	return n
}

// Thread holds a list of nodes and some of its metadata
type Thread struct {
	// ID of the thread
	ThreadID NullInt64 `json:"thread_id,omitempty"`
	// Node from which this thread was forked
	ForkedFrom *Node `json:"forked_from,omitempty"`
	// Team which this node belongs to
	TeamID NullInt64 `json:"team_id,omitempty"`
	// Detail of the source node of this thread
	SourceNode *Node `json:"source_node,omitempty"`
	// Comment nodes of the thread
	ChildNodes []*Node `json:"child_nodes,omitempty"`
	// Whether the thread is open
	IsOpen bool `json:"is_open,omitempty"`
	// Resolved permissions for this thread
	Permissions Permissions `json:"permissions,omitempty"`
	// Extended attributes for the thread
	Attrs JSONB `json:"attrs,omitempty"`
}

// User holds user metadata
type User struct {
	// ID of the user
	UserID NullInt64 `json:"user_id"`
	// Username of the user
	Username NullString `json:"username"`
	// Password of the user
	Password NullString `json:"password"`
	// Primary email of the user
	EmailPrimary NullString `json:"email_primary"`
	// Time when the user account was created
	CreatedAt NullTime `json:"created_at"`
	// Verification status of the user
	IsVerified bool `json:"is_verified"`
	// Active status of the user
	IsActive bool `json:"is_active"`
	// Display name of the user
	DisplayName NullString `json:"display_name,omitempty"`
}

// Team holds team metadata
type Team struct {
	// ID of the team
	TeamID NullInt64 `json:"team_id"`
	// Unique slug of the team
	Slug NullString `json:"slug"`
	// Display name of the team
	DisplayName NullString `json:"display_name,omitempty"`
	// ID of the org the team belongs to
	OrgID NullInt64 `json:"org_id"`
	// ID of the parent team (if nested)
	ParentTeam NullInt64 `json:"parent_team,omitempty"`
	// Resolved permissions for the team
	Permissions Permissions `json:"permissions,omitempty"`
}

// Org holds organization metadata
type Org struct {
	// ID of the organization
	OrgID NullInt64 `json:"org_id"`
	// Unique slug of the organization
	Slug NullString `json:"slug"`
	// Display name of the organization
	DisplayName NullString `json:"display_name"`
	// Primary contact (user) of the organization
	PrimaryContact User `json:"primary_contact"`
}

// AccessLevel holds the type for defining access levels
type AccessLevel string

const (
	// AccessLevelAdmin is used to define admin access
	AccessLevelAdmin = AccessLevel("admin")
	// AccessLevelWrite is used to define read+write access
	AccessLevelWrite = AccessLevel("write")
	// AccessLevelRead is used to define readonly access
	AccessLevelRead = AccessLevel("read")
	// AccessLevelNone is used to define no access
	AccessLevelNone = AccessLevel("none")
)

// Permissions are a matrix of access control levels used by teams and threads
type Permissions struct {
	// Team level permission
	Team AccessLevel `json:"team"`
	// Org level permission
	Org AccessLevel `json:"org"`
	// Public level permission
	Public AccessLevel `json:"public"`
}
