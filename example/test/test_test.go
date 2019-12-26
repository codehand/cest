package test

import "testing"

// TestABC is func test
// Created at 26-12-2019 by tam gh
func TestABC(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ABC()
		})
	}
}
