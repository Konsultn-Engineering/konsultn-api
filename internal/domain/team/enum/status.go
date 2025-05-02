package enum

// Status type
type Status string

// Enum values for roles
const (
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
)

// String method to get string representation of Status
func (s Status) String() string {
	return string(s)
}

// IsValidStatus checks if the given status is valid
func IsValidStatus(status Status) bool {
	switch status {
	case Pending, Accepted, Rejected:
		return true
	}
	return false
}
