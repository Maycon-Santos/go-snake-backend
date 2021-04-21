package uuid

import (
	"github.com/sony/sonyflake"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

func Generate() (*uint64, error) {
	id, err := flake.NextID()
	if err != nil {
		return nil, err
	}

	return &id, err
}
