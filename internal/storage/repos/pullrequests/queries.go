package pullrequests

import _ "embed" // embed package is used to embed SQL query files

var (
	//go:embed queries/create.sql
	createPRQuery string
	//go:embed queries/get_by_id.sql
	getPRByIDQuery string
	//go:embed queries/update.sql
	updatePRQuery string
	//go:embed queries/get_assignments.sql
	getAssignmentsQuery string
	//go:embed queries/get_stats.sql
	getPRStatsQuery string
)
