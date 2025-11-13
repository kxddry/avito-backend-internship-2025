package teams

import _ "embed" // embed package is used to embed SQL query files

var (
	//go:embed queries/create.sql
	createTeamQuery string
	//go:embed queries/get.sql
	getTeamQuery string
	//go:embed queries/get_members.sql
	getTeamMembersQuery string
)
