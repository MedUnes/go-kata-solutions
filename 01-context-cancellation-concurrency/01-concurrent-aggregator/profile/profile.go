package profile

import (
	"context"
)

type Profile struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	Get(ctx context.Context, id int) (*Profile, error)
}
