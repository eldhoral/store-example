package errs

import (
	"database/sql"
	"fmt"
)

// Error Handling, more functions to flexible for Yii2 =================================================================

type ErrorType uint64

const (
	ErrorTypeDefault = iota
	ErrorTypeMissing
	ErrorTypeParse
	ErrorTypeNotFound
	ErrorTypeSql
	ErrorTypeInvalid
	ErrorTypeBadRequest
	ErrorTypeFile
	ErrorTypePanic
	ErrorTypeUnauthorized
	ErrorTypeUnprocesseableEntity
)

const (
	ERROR_400 = "Bad Request"
	ERROR_401 = "Unauthorized"
	ERROR_403 = "Forbidden"
	ERROR_422 = "Unprocessable Entity"
	ERROR_500 = "Internal Server Error"
)

// Error represents a error's specification.
type Error struct {
	Name    string
	Message string
	Err     error
	Type    ErrorType
	File    string
	Stack   string
}

// Error implements the error interface.
func (msg Error) Error() string {
	return msg.Err.Error()
}

func (msg *Error) IsType(flags ErrorType) bool {
	return msg.Type == flags
}

func (msg *Error) Wrap(err error) *Error {
	if err == sql.ErrNoRows {
		return generateError("No records", err, ErrorTypeNotFound)
	}

	// Extend more on demand

	return generateError("Unknown", err, ErrorTypeDefault)
}

// Unwrap returns the wrapped error, to allow interoperability with errors.Is(), errors.As() and errors.Unwrap()
func (msg *Error) Unwrap() error {
	return msg.Err
}

type ErrorDebugResponse struct {
	Message string `json:"message"`
	File    string `json:"file"`
}

// Often in production, we don't show File, and LineNumber for security
func (msg *Error) GetErrorDebugResponse() *ErrorDebugResponse {
	e := ErrorDebugResponse{}
	e.Message = msg.Err.Error()
	e.File = msg.File // No need to show LineNumber, File already include LineNumber

	return &e
}

// Check Type ==========================================================================================================

// ToCustomError short-cut for casting type e, ok := err.(*Error); ok
func ToCustomError(err error) (*Error, bool) {
	if e, ok := err.(*Error); ok {
		return e, ok
	}
	return nil, false
}

func AsType(err error, errType ErrorType) bool {
	if e, ok := err.(*Error); ok {
		if e.Type == errType {
			return true
		}
	}
	return false
}

func AsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		if e.Type == ErrorTypeNotFound {
			return true
		}
	}
	return false
}

func AsNotAuthorized(err error) bool {
	if e, ok := err.(*Error); ok {
		if e.Type == ErrorTypeUnauthorized {
			return true
		}
	}
	return false
}

func AsUnprocessableEntity(err error) bool {
	if e, ok := err.(*Error); ok {
		if e.Type == ErrorTypeUnprocesseableEntity {
			return true
		}
	}
	return false
}

func NotFoundOrSqlError(name string, err error) error {
	if err == sql.ErrNoRows {
		return generateError(name, fmt.Errorf(name), ErrorTypeNotFound)
	} else if err != nil {
		return generateError(name, err, ErrorTypeSql)
	}
	return nil
}

// Shortcut to New Custom Error ========================================================================================

// new Error: is private for correct stack and file
func NewError(name string, e error, errType ErrorType) *Error {
	return generateError(name, e, errType)
}

func generateError(name string, err error, errType ErrorType) *Error {
	stack, file := StackAndFile(3)
	return &Error{
		Message: name,
		Name:    name, // Name for developer to understand and debug
		Err:     err,  // Most response messages get from Err.Error()
		Type:    errType,

		File:  file,
		Stack: stack,
	}
}

func PanicError(name string, err error, skip int) *Error {
	stack, file := StackAndFile(skip)
	return &Error{
		Name:  name,
		Err:   err,
		Type:  ErrorTypePanic,
		File:  file,
		Stack: stack,
	}
}

func NewNotFoundError(name string) *Error {
	return generateError(name, fmt.Errorf(name), ErrorTypeNotFound)
}

func NewSqlError(name string, err error) *Error {
	return generateError(name, err, ErrorTypeSql)
}

func NewBadRequestError(name string, err error) *Error {
	return generateError("Bad request", fmt.Errorf("%w. %s", err, name), ErrorTypeBadRequest)
}

func NewParseError(name string) *Error {
	return generateError(name, fmt.Errorf("Parse :%s fail", name), ErrorTypeParse)
}

func NewMissingQueryError(name string) *Error {
	return generateError(name, fmt.Errorf("Missing query: %s in URI", name), ErrorTypeMissing)
}

func NewMissingKeyError(name string) *Error {
	return generateError(name, fmt.Errorf("Missing key: %s in JSON body", name), ErrorTypeMissing)
}

func NewMissingError(name string) *Error {
	return generateError(name, fmt.Errorf("Missing required parameter: %s", name), ErrorTypeMissing)
}

func NewMalformedJSONError(err error) *Error {
	return generateError("Malformed JSON", err, ErrorTypeParse)
}

func NewMalformedFormError(err error) *Error {
	return generateError("Malformed Form", err, ErrorTypeParse)
}

func NewUnautorizedError(err error) *Error {
	return generateError("Request Unauthorized", err, ErrorTypeUnauthorized)
}

func NewUnprocessableEntityError(err error) *Error {
	return generateError("Unprocessable Entity", err, ErrorTypeUnprocesseableEntity)
}

func NewBadRequest(err error) *Error {
	return &Error{
		Name: "error",
		Err:  err,
		Type: ErrorTypeBadRequest,
	}
}
