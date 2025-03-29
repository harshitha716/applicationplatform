package apierrors

type APIError struct {
	code    int
	message string
	detail  string
}

func (e *APIError) Error() string {
	return e.message
}

func (e *APIError) StatusCode() int {
	return e.code
}

func (e *APIError) Detail() string {
	return e.detail
}
