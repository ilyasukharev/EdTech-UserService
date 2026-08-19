package utils

import (
	"math/rand/v2"
	"strconv"
	"time"
)

func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func StringPtr(s string) *string {
	return &s
}

func GenerateRandNumbersAsString(count int, rightBoundExclude int) string {
	if count <= 0 {
		count = 1
	}

	var str string
	for i := 0; i < count; i++ {
		str += strconv.Itoa(rand.IntN(rightBoundExclude))
	}
	return str
}
