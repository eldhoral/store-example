package realiphelper

import (
	"testing"

	"net/http"

	"github.com/stretchr/testify/assert"
)

func TestIsPrivateAddress(t *testing.T) {
	t.Run("Valid public ip address returns false", func(t *testing.T) {
		example_ip := "192.0.2.1"
		result, err := isPrivateAddress(example_ip)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, false)
	})

	t.Run("Invalid ip address returns false", func(t *testing.T) {
		example_ip := "192.0.2"
		result, err := isPrivateAddress(example_ip)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, false)
	})

	t.Run("Valid private ip address returns true", func(t *testing.T) {
		example_ip := "172.16.0.0"
		result, err := isPrivateAddress(example_ip)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, true)
	})

	t.Run("Empty ip address returns false", func(t *testing.T) {
		example_ip := ""
		result, err := isPrivateAddress(example_ip)

		assert.NotNil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, result, false)
	})
}

func TestFromRequest(t *testing.T) {
	t.Run("Empty headers and empty address returns empty string", func(t *testing.T) {
		mockReq := http.Request{}
		result := FromRequest(&mockReq)

		assert.NotNil(t, result)
		assert.Equal(t, result, "")
	})

	t.Run("Empty headers and host address with port returns host address", func(t *testing.T) {
		mockReq := http.Request{}
		mockReq.RemoteAddr = "localhost:3000"
		result := FromRequest(&mockReq)

		assert.NotNil(t, result)
		assert.Equal(t, result, "localhost")
	})

	t.Run("Empty headers with host address returns host address", func(t *testing.T) {
		mockReq := http.Request{}
		mockReq.RemoteAddr = "127.0.0.1"

		result := FromRequest(&mockReq)

		assert.NotNil(t, result)
		assert.Equal(t, result, "127.0.0.1")
	})

	t.Run("Empty headers with host address of ::1 returns 127.0.0.1", func(t *testing.T) {
		mockReq := http.Request{}
		mockReq.RemoteAddr = "[::1]:80"
		result := FromRequest(&mockReq)

		assert.NotNil(t, result)
		assert.Equal(t, result, "127.0.0.1")
	})

	t.Run("Headers present and a public address in X-Forwarded-For returns public address value", func(t *testing.T) {
		mockReq := http.Request{}
		mockReq.Header = http.Header{}
		mockReq.Header.Add("X-Forwarded-For", "192.0.2.1")
		mockReq.Header.Add("X-Real-Ip", "127.0.0.1")

		result := FromRequest(&mockReq)

		assert.NotNil(t, result)
		assert.Equal(t, result, "192.0.2.1")
	})

	t.Run("Headers present and no public address in X-Forwarded-For returns X-Real-IP address", func(t *testing.T) {
		mockReq := http.Request{}
		mockReq.Header = http.Header{}
		mockReq.Header.Add("X-Forwarded-For", "172.16.0.0")
		mockReq.Header.Add("X-Real-Ip", "127.0.0.1")

		result := FromRequest(&mockReq)

		assert.NotNil(t, result)
		assert.Equal(t, result, "127.0.0.1")
	})
}
