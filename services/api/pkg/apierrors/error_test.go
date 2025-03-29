package apierrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*

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

*/

func TestError(t *testing.T) {
	apiError := &APIError{
		code:    400,
		message: "Bad Request",
		detail:  "msg",
	}

	assert.Equal(t, "Bad Request", apiError.Error())
	assert.Equal(t, 400, apiError.StatusCode())
	assert.Equal(t, "msg", apiError.Detail())

}
