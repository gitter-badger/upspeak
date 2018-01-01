package models

type Team struct {
	ID          int64
	Slug        string
	DisplayName string
	OrgId       string
	ParentTeam  int64
	Permissions JSONB
}
