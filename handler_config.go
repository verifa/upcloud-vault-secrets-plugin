package upcloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func configPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config",

			Fields: map[string]*framework.FieldSchema{
				"username": {
					Type:        framework.TypeString,
					Required:    true,
					Description: "Specifies the upcloud Admin username to authenticate.",
				},
				"password": {
					Type:        framework.TypeString,
					Required:    true,
					Description: "Specifies the upcloud Admin username's password to authenticate.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: handleConfigRead,
					Summary:  "Retrieve the secret from the map.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: handleConfigWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: handleConfigWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: handleConfigDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
			},

			//ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get(ctx, "config")
	if err != nil {
		return nil, fmt.Errorf("could not fetch StorageEntry config: %w", err)
	}
	if entry == nil {
		return logical.ErrorResponse("backend not configured"), nil
	}
	var auth upcloudAuth
	if err := entry.DecodeJSON(&auth); err != nil {
		return nil, fmt.Errorf("could not decode JSON while fetching StorageEntry config: %w", err)
	}
	return &logical.Response{Data: map[string]interface{}{
		"username": auth.Username,
		"password": auth.Password,
	}}, nil
}

func handleConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Check to make sure that kv pairs provided
	if len(req.Data) == 0 {
		return nil, errors.New("data must be provided to store in secret")
	}

	username := data.Get("username")
	if username == nil {
		return nil, errors.New("must provide a username")

	}

	password := data.Get("password")
	if password == nil {
		return nil, errors.New("must provide a password")

	}

	entry, err := logical.StorageEntryJSON("config", upcloudAuth{
		Username: username.(string),
		Password: password.(string),
	})
	if err != nil {
		return nil, fmt.Errorf("could not marshal json for upcloud auth: %w", err)
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, fmt.Errorf("could not put upcloudAuth to storage: %w", err)
	}

	return nil, nil
}

func handleConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "config"); err != nil {
		return nil, fmt.Errorf("failed to delete config: %w", err)
	}

	return nil, nil
}
