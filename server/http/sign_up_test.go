package http

import (
	"reflect"
	"testing"

	"github.com/Maycon-Santos/go-snake-backend/container"
	"github.com/julienschmidt/httprouter"
)

func TestSignUpHandler(t *testing.T) {
	type args struct {
		container container.Container
	}
	tests := []struct {
		name string
		args args
		want httprouter.Handle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SignUpHandler(tt.args.container); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SignUpHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
