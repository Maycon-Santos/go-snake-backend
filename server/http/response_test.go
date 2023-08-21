package http

import (
	"context"
	"net/http"
	"testing"
)

func Test_makeResponse(t *testing.T) {
	type args struct {
		ctx      context.Context
		writer   http.ResponseWriter
		response responseConfig
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
			if err := makeResponse(tt.args.ctx, tt.args.writer, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("makeResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
