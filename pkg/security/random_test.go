package security

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomStringYii(t *testing.T) {
	testCases := []int{32, 0}

	for _, length := range testCases {
		testTitle := fmt.Sprintf("Test Generate length %d characters", length)

		t.Run(testTitle, func(t *testing.T) {
			result := GenerateRandomStringYii(length)

			assert.Equal(t, "string", reflect.TypeOf(result).Name())
			assert.Equal(t, length, len(result))
		})
	}
}
