package upcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const upcloudConfigPath = "config"

// upcloudConfig contains the configuration stored for the upcloud backend,
// including authentication and default values
type upcloudConfig struct {
	// Username for upcloud auth
	Username string `json:"username"`
	// Password for upcloud auth
	Password string `json:"password"`
}

func configPaths(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: upcloudConfigPath,

			Fields: map[string]*framework.FieldSchema{
				"username": {
					Type:        framework.TypeString,
					Description: "Specifies the upcloud Admin username to authenticate.",
				},
				"password": {
					Type:        framework.TypeString,
					Description: "Specifies the upcloud Admin username's password to authenticate.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleConfigRead,
					Summary:  "Retrieve the secret from the map.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleConfigWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleConfigWrite,
					Summary:  "Store a secret at the specified location.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleConfigDelete,
					Summary:  "Deletes the secret at the specified location.",
				},
			},
		},
	}
}

func (b *backend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := pathConfigRead(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("reading upcloud config: %w", err)
	}
	if config == nil {
		return logical.ErrorResponse("config is nil"), nil
	}
	return &logical.Response{Data: map[string]interface{}{
		"username": config.Username,
		"password": config.Password,
	}}, nil
}

func (b *backend) handleConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Read existing config first, if any, so that values can be updated
	config, err := pathConfigRead(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("reading upcloud config: %w", err)
	}
	if config == nil {
		config = new(upcloudConfig)
	}

	username, ok := data.GetOk("username")
	if !ok {
		return logical.ErrorResponse("username must be provided"), nil
	}
	config.Username = username.(string)

	password, ok := data.GetOk("password")
	if !ok {
		return logical.ErrorResponse("password must be provided"), nil
	}
	config.Password = password.(string)

	if err := pathConfigWrite(ctx, req.Storage, *config); err != nil {
		return nil, fmt.Errorf("writing upcloud config: %w", err)
	}

	return nil, nil
}

func (b *backend) handleConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, upcloudConfigPath); err != nil {
		return nil, fmt.Errorf("failed to delete config: %w", err)
	}

	return nil, nil
}

// pathConfigRead reads the upcloud config stored for this backend
func pathConfigRead(ctx context.Context, s logical.Storage) (*upcloudConfig, error) {
	var config upcloudConfig
	configRaw, err := s.Get(ctx, upcloudConfigPath)
	if err != nil {
		return nil, fmt.Errorf("getting upcloud config from storage: %w", err)
	}
	if configRaw == nil {
		return nil, nil
	}
	if err := configRaw.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("decoding upcloud config: %w", err)
	}

	return &config, nil
}

// pathConfigReadMust reads the upcloud config stored for this backend and
// returns an error if no config is set
func pathConfigReadMust(ctx context.Context, s logical.Storage) (*upcloudConfig, error) {
	config, err := pathConfigRead(ctx, s)
	if err != nil {
		// No need to wrap the error on this occassion as we cannot provide
		// any useful context
		return nil, err
	}

	if config == nil {
		return nil, fmt.Errorf("upcloud config must be written to at path <mount>/config")
	}
	return config, nil
}

func pathConfigWrite(ctx context.Context, s logical.Storage, config upcloudConfig) error {
	entry, err := logical.StorageEntryJSON(upcloudConfigPath, config)
	if err != nil {
		return fmt.Errorf("creating storage entry for upcloud config: %w", err)
	}

	if err := s.Put(ctx, entry); err != nil {
		return fmt.Errorf("storing upcloud config: %w", err)
	}
	return nil
}
