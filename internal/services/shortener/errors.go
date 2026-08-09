package shortener

// ValidationError represents a validation error, typically used when input data does not meet certain criteria.
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new ValidationError with the given message.
func NewValidationError(message, field string) *ValidationError {
	return &ValidationError{Message: message, Field: field}
}

// NotFoundError represents a not found error, typically used when a resource cannot be found.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError creates a new NotFoundError with the given message.
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

// ConflictError represents a conflict error, typically used when a resource already exists.
type ConflictError struct {
	Message string
	Field   string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// NewConflictError creates a new ConflictError with the given message.
func NewConflictError(message, field string) *ConflictError {
	return &ConflictError{Message: message, Field: field}
}
