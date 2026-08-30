package errors

// ErrorCode defines the standard SNISID domain error types.
type ErrorCode string

const (
	// InvalidArgument indicates client provided invalid data.
	InvalidArgument ErrorCode = "invalid_argument"

	// NotFound indicates the requested resource does not exist.
	NotFound ErrorCode = "not_found"

	// Conflict indicates a state conflict (e.g., entity already exists).
	Conflict ErrorCode = "conflict"

	// Unauthenticated indicates missing or invalid authentication credentials.
	Unauthenticated ErrorCode = "unauthenticated"

	// PermissionDenied indicates the user does not have permission for the action.
	PermissionDenied ErrorCode = "permission_denied"

	// Internal indicates an unexpected server failure.
	Internal ErrorCode = "internal"

	// Unavailable indicates the service is currently unreachable or down.
	Unavailable ErrorCode = "unavailable"
)
