package rpc

type ThreadService struct {
	Id          int64
	AuthorId    int64
	Permissions struct {
	}
	Attrs struct {
	}
}

// CreateThread creates a thread
func CreateThread() {

}

// UpdateThread updates a thread
func UpdateThread() {

}

// DeleteThread deletes a thread and all its comments
func DeleteThread() {

}

// ArchiveThread archives/unarchives a thread
func ArchiveThread() {

}

// GetThreadComments gets recent 20 comments of a thread unless specified a range
func GetThreadComments() {

}
