package models

type Team struct {
	Id          int64
	Slug        string
	DisplayName string
	OrgId       string
	ParentTeam  int64
	Permissions JSONB
}
