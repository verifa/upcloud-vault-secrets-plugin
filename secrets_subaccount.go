package upcloud

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const secretSubAccountType = "subaccount"

func secretSubAccount(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: secretSubAccountType,
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Required: true,
				Type:     framework.TypeString,
			},
		},
		Revoke: b.secretSubAccountRevoke,
	}
}

func (b *backend) secretSubAccountRevoke(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	username := data.Get("username")
	if username == "" {
		return logical.ErrorResponse("cannot revoke subaccount with an empty username"), nil
	}

	// Uncomment this for access to config
	// config, err := pathConfigReadMust(ctx, req.Storage)
	// if err != nil {
	// 	return nil, fmt.Errorf("reading upcloud config: %w", err)
	// }

	return logical.ErrorResponse("Not implemented"), nil
}
