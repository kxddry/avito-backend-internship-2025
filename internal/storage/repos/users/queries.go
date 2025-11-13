package users

import _ "embed" // embed package is used to embed SQL query files

var (
	//go:embed queries/get_by_id.sql
	getByIDQuery string
	//go:embed queries/update.sql
	updateUserQuery string
	//go:embed queries/upsert.sql
	upsertUserQuery string
)
