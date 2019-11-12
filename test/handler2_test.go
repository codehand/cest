// ABC is func1
package test

import "testing"

func TestABC(t *testing.T) {
	type args struct {
		in int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := ABC(tt.args.in); (err != nil) != tt.wantErr {
			t.Errorf("%q. ABC() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
