package enum

// Role type
type Role string

// Enum values for roles
const (
	Owner  Role = "owner"
	Admin  Role = "admin"
	Member Role = "member"
)

// String method to get string representation of Role
func (r Role) String() string {
	return string(r)
}

// IsValidRole checks if the given role is valid
func IsValidRole(role Role) bool {
	switch role {
	case Owner, Admin, Member:
		return true
	}
	return false
}
