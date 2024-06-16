package slicehelper

import (
	"reflect"
	"strings"
)

//UniqueSlices Remove duplicate elements on slice of string
func UniqueSlices(arr []string) []string {
	size := len(arr)
	result := make([]string, 0, size)
	temp := map[string]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; !ok {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringContainSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

// This function is the same as in_array in PHP
func InArray(key interface{}, array interface{}) bool {
	switch key := key.(type) {
	case string:
		for _, item := range array.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range array.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range array.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func KeyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}

func NilFields(x interface{}) bool {
	rv := reflect.ValueOf(x)
	rv = rv.Elem()

	for i := 0; i < rv.NumField(); i++ {
		if !rv.Field(i).IsNil() {
			return false
		}
	}
	return true
}
