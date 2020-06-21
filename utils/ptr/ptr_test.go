package ptr

import (
	"testing"
)

// TestInt64Ptr verifies Int64Ptr.
func TestInt64Ptr(t *testing.T) {
	n := int64(1000000007)
	if got, want := *Int64Ptr(n), n; got != want {
		t.Errorf("*Int64Ptr(%v) = %v, want %v", n, got, want)
	}
}

// TestInt64 verifies Int64.
func TestInt64(t *testing.T) {
	tests := []struct {
		desc string
		p    *int64
		want int64
	}{
		{
			desc: "nil pointer is converted to 0",
		},
		{
			desc: "non nil pointer",
			p:    Int64Ptr(1000000009),
			want: 1000000009,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got, want := Int64(tt.p), tt.want; got != want {
				t.Errorf("Int64() = %v, want %v", got, want)
			}
		})
	}
}

// TestInt32Ptr verifies Int32Ptr.
func TestInt32Ptr(t *testing.T) {
	n := int32(1000000007)
	if got, want := *Int32Ptr(n), n; got != want {
		t.Errorf("*Int32Ptr(%v) = %v, want %v", n, got, want)
	}
}

// TestBoolPtr verifies BoolPtr.
func TestBoolPtr(t *testing.T) {
	tests := []struct {
		desc string
		b    bool
	}{
		{
			desc: "Get true",
			b:    true,
		},
		{
			desc: "Get false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got, want := *BoolPtr(tt.b), tt.b; got != want {
				t.Errorf("*BoolPtr(%v) = %v, want %v", tt.b, got, want)
			}
		})
	}
}

// TestBool verifies Bool.
func TestBool(t *testing.T) {
	tests := []struct {
		desc string
		p    *bool
		want bool
	}{
		{
			desc: "nil pointer is converted to false",
		},
		{
			desc: "true",
			p:    BoolPtr(true),
			want: true,
		},
		{
			desc: "false",
			p:    BoolPtr(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got, want := Bool(tt.p), tt.want; got != want {
				t.Errorf("Bool(...) = %v, want %v", got, want)
			}
		})
	}
}

// TestStringPtr verifies StringPtr.
func TestStringPtr(t *testing.T) {
	tests := []struct {
		desc string
		s    string
	}{
		{
			desc: "Get string",
			s:    "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got, want := *StringPtr(tt.s), tt.s; got != want {
				t.Errorf("*StringPtr(%q) = %v, want %v", tt.s, got, want)
			}
		})
	}
}

// TestString verifies String.
func TestString(t *testing.T) {
	tests := []struct {
		desc string
		p    *string
		want string
	}{
		{
			desc: "nil pointer is converted to an empty string",
		},
		{
			desc: "non empty string",
			p:    StringPtr("hello world"),
			want: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if got, want := String(tt.p), tt.want; got != want {
				t.Errorf("String(...) = %v, want %v", got, want)
			}
		})
	}
}
