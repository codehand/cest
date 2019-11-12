package test

import "testing"

// ABC is func1

// TestABC is func test
// abcd
func TestABC(t *testing.T) {
	type args struct {
		in int
		x  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases 11.
	}
	for _, tt := range tests {
		if err := ABC(tt.args.in); (err != nil) != tt.wantErr {
			t.Errorf("%q. ABC() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
