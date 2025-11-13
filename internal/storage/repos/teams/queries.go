package teams

import _ "embed"

var (
	//go:embed queries/create.sql
	createTeamQuery string
	//go:embed queries/get.sql
	getTeamQuery string
	//go:embed queries/get_members.sql
	getTeamMembersQuery string
)
