package cconstant

import "time"

const (
	AuthDB     string = "taskdb.public.auth"
	CommandsDB string = "taskdb.public.cmds"
	ResultDB   string = "taskdb.public.results"
)

const (
	Salt            = "xjifcmefdx2oxe3x"
	SignedKey       = "efcj34s3dr4cwdxxjuu34"
	SecterKey       = "top_secret"
	AccessTokenTTL  = 2 * time.Hour
	RefreshTokenTTl = 30 * 24 * time.Hour
)

const (
	RoleUser  int = 0
	RoleAdmin int = 1
)

const (
	RunningStatus = iota
	FailedStatus
	SuccessStatus
)
