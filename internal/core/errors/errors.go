package errors

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type NoAccessError struct {
	Message string
}

func (e *NoAccessError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type MappingError struct {
	Message string
}

func (e *MappingError) Error() string {
	return e.Message
}

type BodyMappingError struct {
	Message string
}

func (e *BodyMappingError) Error() string {
	return e.Message
}
