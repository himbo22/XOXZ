package _const

type Permission string

const (
	ACCOUNT_PROVIDE Permission = "account:provide"
	ACCOUNT_CREATE  Permission = "account:create"
	ACCOUNT_UPDATE  Permission = "account:update"
)
