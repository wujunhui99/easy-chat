package constants

const (
	SYSTEM_ROOT_UID = "root"
)

const (
	UserStatusNormal   int64 = 0
	UserStatusDisabled int64 = 1 << 0
	UserStatusPending  int64 = 1 << 1
	UserStatusDeleted  int64 = 1 << 2
)

const (
	UserTypeHuman int64 = 0
	UserTypeAgent int64 = 1
)

func UserStatusHas(status, flag int64) bool {
	return status&flag != 0
}

func IsUserActive(status int64) bool {
	return !UserStatusHas(status, UserStatusDisabled|UserStatusDeleted)
}

func WithUserStatus(status, flag int64) int64 {
	return status | flag
}

func ClearUserStatus(status, flag int64) int64 {
	return status &^ flag
}
