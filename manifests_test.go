package main

import "testing"

func Test_applyManifests(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := applyManifests(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("applyManifests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
