package account

type AccountRole string

const (
	SuperAdmin  AccountRole = "SUPERADMIN"
	Admin AccountRole = "ADMIN"
	Guest AccountRole = "GUEST"
)