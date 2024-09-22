package converter

const (
	UNKNOWN = 0
	ADMIN   = 1
	USER    = 2
)

func ToRoleString(role int32) string {
	switch role {
	case ADMIN:
		return "admin"
	case USER:
		return "user"
	default:
		return "unknown"
	}
}

func ToRole(role string) int32 {
	switch role {
	case "admin":
		return ADMIN
	case "user":
		return USER
	default:
		return UNKNOWN
	}
}
