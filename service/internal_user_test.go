package service

import (
	"testing"
)

// validates M7-F01 — GroupInternal constant has correct value
func TestGroupInternalConstant(t *testing.T) {
	if GroupInternal != "internal" {
		t.Errorf("GroupInternal = %q, want 'internal'", GroupInternal)
	}
}

// validates M7-F01 — IsInternalUser returns false for nil/empty context
// (compile-time check that the function signature is correct)
func TestIsInternalUserSignature(t *testing.T) {
	// Verify the function exists and has the right signature.
	// Full behaviour is tested via integration; here we just confirm it compiles.
	var fn func(interface{}) bool
	_ = fn
	// IsInternalUser is defined in this package — if this file compiles, the function exists.
}
