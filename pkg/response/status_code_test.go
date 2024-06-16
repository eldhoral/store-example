package response

import (
    "testing"

    "store-api/pkg/errs"

    "github.com/stretchr/testify/assert"
)

func TestParams_GetStatusCode(t *testing.T) {
    t.Run("no-error-200", func(t *testing.T) {
        err := GetStatusCode(nil)
        assert.Equal(t, 200, err)
    })

    t.Run("error-500", func(t *testing.T) {
        err := GetStatusCode(errs.ErrInternalServerError)
        assert.Equal(t, 500, err)
    })

    t.Run("error-404", func(t *testing.T) {
        err := GetStatusCode(errs.ErrNotFound)
        assert.Equal(t, 404, err)
    })

    t.Run("error-409", func(t *testing.T) {
        err := GetStatusCode(errs.ErrConflict)
        assert.Equal(t, 409, err)
    })

    t.Run("error-400", func(t *testing.T) {
        err := GetStatusCode(errs.ErrBadParamInput)
        assert.Equal(t, 400, err)
    })

    t.Run("error-500", func(t *testing.T) {
        err := GetStatusCode(errs.NewCustomError("new-custom-error"))
        assert.Equal(t, 500, err)
    })

    t.Run("error-422", func(t *testing.T) {
        err := GetStatusCode(errs.ErrMissingParam)
        assert.Equal(t, 422, err)
    })
}
