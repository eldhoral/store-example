package converterhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolToInt(t *testing.T) {
	t.Run("true returns 1", func(t *testing.T) {
		result := BoolToInt(true)

		assert.NotNil(t, result)
		assert.Equal(t, result, 1)
	})

	t.Run("false returns 0", func(t *testing.T) {
		result := BoolToInt(false)

		assert.NotNil(t, result)
		assert.Equal(t, result, 0)
	})
}

func TestInterfaceToString(t *testing.T) {
	t.Run("String interface returns non empty string", func(t *testing.T) {
		var example interface{} = "example_string"
		result := InterfaceToString(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "example_string")
	})

	t.Run("Int interface returns non empty string", func(t *testing.T) {
		var example interface{} = 12345
		result := InterfaceToString(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "12345")
	})

	t.Run("Float interface returns non empty string", func(t *testing.T) {
		var example interface{} = 12345.12345
		result := InterfaceToString(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "12345.12345")
	})

	t.Run("Bool interface returns non empty string", func(t *testing.T) {
		var example interface{} = true
		result := InterfaceToString(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "true")
	})

	t.Run("Array interface returns non empty string", func(t *testing.T) {
		var example interface{} = []string{"string1", "string2"}
		result := InterfaceToString(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "[string1 string2]")
	})

	t.Run("Empty interface returns non empty string", func(t *testing.T) {
		var example interface{}

		result := InterfaceToString(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "<nil>")
	})
}

func TestInterfaceToInt64(t *testing.T) {
	t.Run("Positive int64 returns non zero int64", func(t *testing.T) {
		var example int64 = 12345
		result := InterfaceToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(12345))
	})

	t.Run("Negative int64 returns non zero int64", func(t *testing.T) {
		var example int64 = -12345
		result := InterfaceToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(-12345))
	})

	t.Run("Positive int returns non zero int64", func(t *testing.T) {
		var example int = 12345
		result := InterfaceToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(12345))
	})

	t.Run("Negative int returns non zero int64", func(t *testing.T) {
		var example int = -12345
		result := InterfaceToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(-12345))
	})

	t.Run("Float int returns zero value int64", func(t *testing.T) {
		var example float64 = 1234.5678
		result := InterfaceToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(0))
	})

	t.Run("Valid int string returns non zero int64", func(t *testing.T) {
		var example string = "1234567"

		result := InterfaceToInt64(example)
		assert.NotNil(t, result)
		assert.Equal(t, result, int64(1234567))
	})

	t.Run("Invalid string returns zero value int64", func(t *testing.T) {
		var example string = "abcdEFGH"
		result := InterfaceToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(0))
	})

}

func TestInterfaceToInt(t *testing.T) {
	t.Run("Positive int returns non zero int", func(t *testing.T) {
		var example interface{} = 12345
		result := InterfaceToInt(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int(12345))
	})

	t.Run("Negative int returns non zero int", func(t *testing.T) {
		var example interface{} = -12345
		result := InterfaceToInt(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int(-12345))
	})

	t.Run("Positive int string returns non zero int", func(t *testing.T) {
		var example interface{} = "12345"
		result := InterfaceToInt(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int(12345))
	})

	t.Run("Negative int string returns non zero int", func(t *testing.T) {
		var example interface{} = "-12345"
		result := InterfaceToInt(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int(-12345))
	})

	t.Run("Invalid string returns zero value int", func(t *testing.T) {
		var example interface{} = ""
		result := InterfaceToInt(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int(0))
	})
}

func TestIntToBool(t *testing.T) {
	t.Run("1 returns true", func(t *testing.T) {
		var example int = 1
		result := IntToBool(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, true)
	})

	t.Run("0 returns false", func(t *testing.T) {
		var example int = 0
		result := IntToBool(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, false)
	})

	t.Run("Number larger than 1 returns false", func(t *testing.T) {
		var example int = 432
		result := IntToBool(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, false)
	})

	t.Run("Number less than 0 returns false", func(t *testing.T) {
		var example int = -5001
		result := IntToBool(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, false)
	})
}

func TestStringToInt64(t *testing.T) {
	t.Run("Positive int string returns non zero int64", func(t *testing.T) {
		var example string = "1231231312"
		result := StringToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(1231231312))
	})

	t.Run("Negative int string returns non zero int64", func(t *testing.T) {
		var example string = "-4563"
		result := StringToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(-4563))
	})

	t.Run("Float string returns zero value int64", func(t *testing.T) {
		var example string = "1222.019"
		result := StringToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(0))
	})

	t.Run("Alphabetical string returns zero value int64", func(t *testing.T) {
		var example string = "abbjjeKLOMOP"
		result := StringToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(0))
	})

	t.Run("Empty string returns zero value int64", func(t *testing.T) {
		var example string = ""
		result := StringToInt64(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, int64(0))
	})
}

func TestUcfirst(t *testing.T) {
	t.Run("Empty string returns empty string", func(t *testing.T) {
		var example string = ""
		result := Ucfirst(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "")
	})

	t.Run("All lowercase string returns all lowercase string (excluding first letter)", func(t *testing.T) {
		var example string = "thisisalowercasestring abc"
		result := Ucfirst(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "Thisisalowercasestring abc")
	})

	t.Run("All capitalized string returns all capitalized string.", func(t *testing.T) {
		var example string = "THISISASTRING ABCDE"
		result := Ucfirst(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "THISISASTRING ABCDE")
	})

	t.Run("Number string returns number string", func(t *testing.T) {
		var example string = "90012"
		result := Ucfirst(example)

		assert.NotNil(t, result)
		assert.Equal(t, result, "90012")
	})
}
