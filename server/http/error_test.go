package http

import (
	"context"
	"testing"
)

func Test_handleError(t *testing.T) {
	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleError(tt.args.ctx, tt.args.err)
		})
	}
}
