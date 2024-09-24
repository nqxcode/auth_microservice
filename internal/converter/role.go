package converter

const (
	unknown = "unknown"
	admin   = "admin"
	user    = "user"
)

// ToRoleString to role string
func ToRoleString(role int32) string {
	switch role {
	case 1:
		return admin
	case 2:
		return user
	default:
		return unknown
	}
}

// ToRole to role
func ToRole(role string) int32 {
	switch role {
	case admin:
		return 1
	case user:
		return 2
	default:
		return 0
	}
}
