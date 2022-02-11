package upcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// backend wraps the backend framework and adds a map for storing key value pairs
type backend struct {
	*framework.Backend

	store map[string][]byte
}

type upcloudAuth struct {
	// Admin username for upcloud auth
	Username string `json:"username"`

	// Admin password for upcloud auth
	Password string `json:"password"`
}

var _ logical.Factory = Factory

// Factory configures and returns Mock backends
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := newBackend()
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, fmt.Errorf("configuration passed into backend is nil")
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func newBackend() (*backend, error) {
	b := &backend{
		store: make(map[string][]byte),
	}

	b.Backend = &framework.Backend{
		Help:        strings.TrimSpace(mockHelp),
		BackendType: logical.TypeLogical,
		Paths: framework.PathAppend(
			configPaths(),
			b.subaccountPaths(),
		),
		Secrets: []*framework.Secret{
			{
				Type: "mike",
				Fields: map[string]*framework.FieldSchema{
					"beer": {
						Required: true,
						Type:     framework.TypeString,
					},
				},
				Revoke: func(c context.Context, r *logical.Request, fd *framework.FieldData) (*logical.Response, error) {
					return logical.ErrorResponse("mike was 'ere: %v", fd), fmt.Errorf("mike was 'ere: %v", fd)
				},
			},
		},
	}

	return b, nil
}

const mockHelp = `
The Mock backend is a dummy secrets backend that stores kv pairs in a map.
`
