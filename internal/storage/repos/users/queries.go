package users

import _ "embed"

var (
	//go:embed queries/get_by_id.sql
	getByIDQuery string
	//go:embed queries/update.sql
	updateUserQuery string
	//go:embed queries/upsert.sql
	upsertUserQuery string
)
