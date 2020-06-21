package ptr

// Int64Ptr returns the pointer to int64 of value n.
func Int64Ptr(n int64) *int64 {
	return &n
}

// Int64 returns the value of *p if p is not nil.  Otherwise, it returns 0.
func Int64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// Int32Ptr returns the pointer to int32 of value n.
func Int32Ptr(n int32) *int32 {
	return &n
}

// BoolPtr returns the pointer to bool of value b.
func BoolPtr(b bool) *bool {
	return &b
}

// Bool returns the value of *p if p is not nil.  Otherwise, it returns false.
func Bool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// StringPtr returns the pointer to string of value s.
func StringPtr(s string) *string {
	return &s
}

// String returns the value of *p if p is not nil.  Otherwise, it returns an empty string.
func String(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
