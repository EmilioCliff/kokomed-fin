package pkg

import (
	"strconv"
	"strings"
	"time"
)

func StringToUint32(s string) (uint32, error) {
	// if s == "" {
	// 	return 0, Errorf(INVALID_ERROR, "id/page is required")
	// }
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, Errorf(INVALID_ERROR, "invalid id/page: %s", err.Error())
	}

	return uint32(id), nil
}

func StringToFloat64(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, Errorf(INVALID_ERROR, "input string is empty")
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, Errorf(INVALID_ERROR, "invalid float64 format: "+s)
	}

	return f, nil
}

func PtrToStr(s *string) string { return *s }

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string { return &s }

// Uint32Ptr returns a pointer to the given uint32.
func Uint32Ptr(i uint32) *uint32 { return &i }

// Float64Ptr returns a pointer to the given float64.
func Float64Ptr(f float64) *float64 { return &f }

// BoolPtr returns a pointer to the given bool.
func BoolPtr(b bool) *bool { return &b }

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int { return &i }

// Int32Ptr returns a pointer to the given int32.
func Int32Ptr(i int32) *int32 { return &i }

// Int64Ptr returns a pointer to the given int64.
func Int64Ptr(i int64) *int64 { return &i }

// UintPtr returns a pointer to the given uint.
func UintPtr(i uint) *uint { return &i }

// Uint64Ptr returns a pointer to the given uint64.
func Uint64Ptr(i uint64) *uint64 { return &i }

// Float32Ptr returns a pointer to the given float32.
func Float32Ptr(f float32) *float32 { return &f }

// BytePtr returns a pointer to the given byte.
func BytePtr(b byte) *byte { return &b }

// RunePtr returns a pointer to the given rune.
func RunePtr(r rune) *rune { return &r }

// TimePtr returns a pointer to the given time.Time.
func TimePtr(t time.Time) *time.Time { return &t }

// DurationPtr returns a pointer to the given time.Duration.
func DurationPtr(d time.Duration) *time.Duration { return &d }

// transform interface to float64 if error it logs and return 0.00
// used to trransform money from db to go type float64
func InterfaceFloat64(i interface{}) float64 {
	iByte, ok := i.([]byte)
	if !ok {
		return 0.00
	}
	iFloat, err := strconv.ParseFloat(string(iByte), 64)
	if err != nil {
		return 0.00
	}

	return iFloat
}
