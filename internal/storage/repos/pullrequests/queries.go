package pullrequests

import _ "embed"

var (
	//go:embed queries/create.sql
	createPRQuery string
	//go:embed queries/get_by_id.sql
	getPRByIDQuery string
	//go:embed queries/update.sql
	updatePRQuery string
	//go:embed queries/get_assignments.sql
	getAssignmentsQuery string
)
