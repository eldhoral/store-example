package response

import (
    "net/http"

    "store-api/pkg/errs"
)

// GetStatusCode represent http status code
func GetStatusCode(err error) int {
    if err == nil {
        return http.StatusOK
    }

    switch err {
    case errs.ErrInternalServerError:
        return http.StatusInternalServerError
    case errs.ErrNotFound:
        return http.StatusNotFound
    case errs.ErrConflict:
        return http.StatusConflict
    case errs.ErrBadParamInput:
        return http.StatusBadRequest
    case errs.ErrMissingParam:
        return http.StatusUnprocessableEntity
    default:
        return http.StatusInternalServerError
    }
}
