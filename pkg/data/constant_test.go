package data

import (
    "testing"

    "store-api/pkg/data/constant"

    "github.com/stretchr/testify/assert"
)

func TestConst(t *testing.T) {

    // DefaultTimeLayout must equal 2006-01-02 15:04:05
    assert.Equal(t, constant.DefaultDatetimeLayout, "2006-01-02 15:04:05")
    assert.Equal(t, constant.DefaultDateLayout, "2006-01-02")
    assert.Equal(t, constant.DefaultTimeLayout, "15:04:05")
}
