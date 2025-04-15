package dashboard

import (
	"context"
)

//after writing the function we should not give the space in between the api and finction ..

// encore:api public path=/hi/:name
func Welcome(ctx context.Context, name string) (*Response, error) {
	msg := "Hello " + name + ", welcome!"
	return &Response{Message: msg}, nil
}

type Response struct {
	Message string
}
