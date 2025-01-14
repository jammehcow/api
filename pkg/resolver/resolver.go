package resolver

import (
	"context"
	"net/http"
	"net/url"
)

type Resolver interface {
	Check(ctx context.Context, url *url.URL) (context.Context, bool)
	Run(ctx context.Context, url *url.URL, r *http.Request) ([]byte, error)
	Name() string
}
