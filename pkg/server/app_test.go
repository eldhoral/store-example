package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOk(t *testing.T) {
	response := Response{
		Status: 200,
	}

	t.Run("Success Read response status", func(t *testing.T) {
		assert.Equal(t, true, response.IsOk())
	})
}

func TestGetDataAsString(t *testing.T) {
	response := Response{
		Data: "Testing",
	}

	t.Run("Success Read string data", func(t *testing.T) {
		assert.Equal(t, response.Data, response.GetDataAsString())
	})

	t.Run("Return empty string if data is not string type", func(t *testing.T) {
		response.Data = 123
		assert.Equal(t, "", response.GetDataAsString())
	})
}

func TestGetStatus(t *testing.T) {
	response := Response{
		Status: 401,
	}

	t.Run("Success read status", func(t *testing.T) {
		assert.Equal(t, response.Status, response.GetStatus())
	})
}

func TestGetYii2Exception(t *testing.T) {
	data := NotAllowedMethod{
		Name: "My error",
	}
	response := Response{
		Data: data,
	}

	t.Run("Success read data with no error", func(t *testing.T) {
		myData, ok := response.GetMobileException().(NotAllowedMethod)

		assert.Equal(t, true, ok)
		assert.Equal(t, myData, data)
	})
}

func TestGetStream(t *testing.T) {
	var myByte bytes.Buffer

	_, err := myByte.Write([]byte{'h', 'i'})
	if err != nil {
		return
	}
	response := Response{
		Data: myByte,
	}

	t.Run("Success read data as bytes", func(t *testing.T) {
		asBuffer := response.Data.(bytes.Buffer)
		assert.Equal(t, asBuffer, response.GetStream())
	})
}
