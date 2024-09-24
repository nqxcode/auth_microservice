package converter

const (
	UNKNOWN = "unknown"
	ADMIN   = "admin"
	USER    = "user"
)

func ToRoleString(role int32) string {
	switch role {
	case 1:
		return ADMIN
	case 2:
		return USER
	default:
		return UNKNOWN
	}
}

func ToRole(role string) int32 {
	switch role {
	case ADMIN:
		return 1
	case USER:
		return 2
	default:
		return 0
	}
}
