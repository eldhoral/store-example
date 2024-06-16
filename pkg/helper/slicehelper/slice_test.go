package slicehelper_test

import (
    "testing"

    "store-api/pkg/helper/slicehelper"

    "github.com/stretchr/testify/assert"
)

func TestUniqueSlices(t *testing.T) {
    t.Run("Return unique slices", func(t *testing.T) {
        duplicateArr := []string{"i", "u", "u"}
        uniqueArr := []string{"i", "u"}

        assert.Equal(t, uniqueArr, slicehelper.UniqueSlices(duplicateArr))
    })

    t.Run("Return same array if no duplicate", func(t *testing.T) {
        inputArr := []string{"i", "u"}
        uniqueArr := []string{"i", "u"}

        assert.Equal(t, uniqueArr, slicehelper.UniqueSlices(inputArr))
    })

    t.Run("Return empty array", func(t *testing.T) {
        inputArr := []string{}

        assert.Equal(t, inputArr, slicehelper.UniqueSlices(inputArr))
    })
}

func TestStringInSlice(t *testing.T) {
    myArr := []string{"i", "u", "aaaaa"}

    t.Run("Return true if item is existed", func(t *testing.T) {
        find := "i"

        assert.Equal(t, true, slicehelper.StringInSlice(find, myArr))
    })

    t.Run("Return false if not existed", func(t *testing.T) {
        find := "x"

        assert.Equal(t, false, slicehelper.StringInSlice(find, myArr))
    })
}

func TestInArray(t *testing.T) {
    myArr := []string{"i", "u", "aaaaa"}

    t.Run("Return true if item is existed - string slice", func(t *testing.T) {
        find := "i"

        assert.Equal(t, true, slicehelper.InArray(find, myArr))
    })

    t.Run("Return true if exist - int slice", func(t *testing.T) {
        find := 11
        myArr := []int{11, 22}

        assert.Equal(t, true, slicehelper.InArray(find, myArr))
    })

    t.Run("Return true if exist - int64 slice", func(t *testing.T) {
        find := int64(22)
        myArr := []int64{11, 22}

        assert.Equal(t, true, slicehelper.InArray(find, myArr))
    })

    t.Run("Return false if not existed", func(t *testing.T) {
        find := "x"

        assert.Equal(t, false, slicehelper.InArray(find, myArr))
    })

    t.Run("Return false if not valid supported type", func(t *testing.T) {
        find := true // boolean key is not supported in InArray function

        assert.Equal(t, false, slicehelper.InArray(find, myArr))
    })
}

func TestKeyExists(t *testing.T) {
    myMap := map[string]interface{}{
        "my":  123,
        "one": 111,
        "nil": nil,
    }

    t.Run("Return true if key exist", func(t *testing.T) {
        find := "my"

        assert.Equal(t, true, slicehelper.KeyExists(myMap, find))
    })

    t.Run("Return false if key not existed", func(t *testing.T) {
        find := "x"

        assert.Equal(t, false, slicehelper.KeyExists(myMap, find))
    })

    t.Run("Return false if key exist but nil value", func(t *testing.T) {
        find := "nil"

        assert.Equal(t, false, slicehelper.KeyExists(myMap, find))
    })
}
