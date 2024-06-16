package data

import (
    "testing"
    "time"

    "store-api/pkg/helper/timehelper"

    "github.com/stretchr/testify/assert"
)

func TestParams_GetInt(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("offset", 0)
    params.Add("limit", "1")
    params.Add("sso_id", "abc-def")

    assert.Equal(t, 0, params.GetInt("offset"))
    assert.Equal(t, 1, params.GetInt("limit"))
    assert.Equal(t, 0, params.GetInt("sso_id"))
}

func TestParams_Add(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("offset", 0)

    assert.Equal(t, params, params)
    assert.NotEqual(t, "a", params)
}

func TestParams_GetInt64(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("offset", int64(0))
    params.Add("limit", 10)

    assert.Equal(t, int64(0), params.GetInt64("offset"))
    assert.Equal(t, int64(10), params.GetInt64("limit"))
}

func TestParams_GetTime(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("due_on", timehelper.NowInUTCAsDateTimeStr())
    x, err := params.GetTime("due_on")

    assert.Equal(t, timehelper.ToLocal(time.Now()).Day(), x.Local().Day())
    assert.NoError(t, err)
}

func TestParams_GetLocalTime(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("due_on", timehelper.NowInUTCAsDateTimeStr())
    x, err := params.GetLocalTime("due_on")

    assert.Equal(t, timehelper.ToLocal(time.Now()).Day(), x.Local().Day())
    assert.NoError(t, err)
}

func TestParams_GetValue(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("offset", 0)
    params.GetValue("offset")

    assert.Equal(t, 0, params.GetValue("offset"))
    assert.NotEqual(t, 1, params.GetValue("offset"))
}

func TestParams_GetString(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("ssoid", "abc-def")
    params.Add("offset", 0)

    assert.Equal(t, "abc-def", params.GetString("ssoid"))
    assert.NotEqual(t, 1, params.GetString("ssoid"))
    assert.Equal(t, "", params.GetString("offset"))
}

func TestParams_GetBool(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("IsNoDueDate", true)

    assert.Equal(t, true, params.GetBool("IsNoDueDate"))
    assert.NotEqual(t, false, params.GetValue("offset"))
}

func TestParams_GetIntOrError(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("offset", 1)
    params.Add("limit", nil)
    x, _ := params.GetIntOrError("offset")
    y, err := params.GetIntOrError("limit")

    assert.Equal(t, 1, x)
    assert.Equal(t, 0, y)
    assert.EqualError(t, err, "Key not found for "+"limit")
}

func TestParams_HasKey(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("haskey", "abcdef")
    params.Add("haskeyUser", nil)

    assert.Equal(t, false, params.HasKey("hashkey"))
    assert.NotEqual(t, true, params.HasKey("haskeyUser"))
}

func TestParams_ConvertInterfaceToInt(t *testing.T) {
    params := NewParamsWrapper()
    //convert from string
    params.Add("offset", "1")
    params.Add("limit", int64(1))
    params.Add("Pph21Paid", float64(1))
    params.Add("cols", []string{"1"})
    w, err := params.ConvertInterfaceToInt("cols")
    x, err := params.ConvertInterfaceToInt("offset")
    y, err := params.ConvertInterfaceToInt("limit")
    z, err := params.ConvertInterfaceToInt("Pph21Paid")

    //convert err return 0
    t.Run("error-return-0", func(t *testing.T) {
        assert.Equal(t, 0, w)
    })

    //convert from string
    t.Run("success-from-string", func(t *testing.T) {
        assert.Equal(t, 1, x)
        assert.Equal(t, true, err)
    })

    //convert from int64
    t.Run("success-from-int64", func(t *testing.T) {
        assert.Equal(t, 1, y)
        assert.Equal(t, true, err)
    })

    //convert from float64
    t.Run("success-from-float64", func(t *testing.T) {
        assert.Equal(t, 1, z)
        assert.Equal(t, true, err)
    })
}

func TestParams_ConvertInterfaceToInt64(t *testing.T) {
    params := NewParamsWrapper()
    params.Add("offset", "1")
    params.Add("limit", 1)
    params.Add("Pph21Paid", float64(1))
    params.Add("cols", []string{"1"})
    w, err := params.ConvertInterfaceToInt64("cols")
    x, err := params.ConvertInterfaceToInt64("offset")
    y, err := params.ConvertInterfaceToInt64("limit")
    z, err := params.ConvertInterfaceToInt64("Pph21Paid")

    //convert err return 0
    t.Run("error-return-0", func(t *testing.T) {
        assert.Equal(t, int64(0), w)
    })

    //convert from string
    t.Run("success-from-string", func(t *testing.T) {
        assert.Equal(t, int64(1), x)
        assert.Equal(t, true, err)
    })

    //convert from int
    t.Run("success-from-int", func(t *testing.T) {
        assert.Equal(t, int64(1), y)
        assert.Equal(t, true, err)
    })

    //convert from float64
    t.Run("success-from-float64", func(t *testing.T) {
        assert.Equal(t, int64(1), z)
        assert.Equal(t, true, err)
    })
}
