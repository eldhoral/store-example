package errs

import (
    "errors"
    "fmt"
    "strings"

    "store-api/pkg/data/filedata"
)

var (
    //ErrNoSuchFile indicates no file found
    ErrNoSuchFile = "open .env: no such file or directory"
    // ErrInternalServerError will throw if any the Internal Server Error happen
    ErrInternalServerError = errors.New("internal server error")
    // ErrNotFound will throw if the requested data is not exists
    ErrNotFound = errors.New("your requested data is not found")
    // ErrConflict will throw if the current action already exists
    ErrConflict = errors.New("your data already exists")
    // ErrBadParamInput will throw if the given request-body or params is not valid
    ErrBadParamInput = errors.New("the given parameter is invalid")
    // ErrMissingParam will throw if the some given request params is missing
    ErrMissingParam = errors.New("the given parameter is missing")
    // ErrFileNotSupported will throw if uploaded file types are not supported in file util
    ErrFileNotSupported = fmt.Errorf("file extension must be between the following types: %s", strings.Join(filedata.AllowedFile, ", "))
    // Color format is invalid, i.e. must be RGB or RGBA
    ErrInvalidColor = errors.New("color is invalid")
)
