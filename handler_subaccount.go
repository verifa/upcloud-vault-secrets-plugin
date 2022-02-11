package upcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type upcloudSubaccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AllowAPI int    `json:"allow_api"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func subaccountPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "subaccount/" + framework.GenericNameRegex("subaccount"),

			Fields: map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				// logical.ReadOperation: &framework.PathOperation{
				// 	Callback: handleSubaccountRead,
				// 	Summary:  "Retrieve the UpCloud subaccount credentials.",
				// },
				logical.UpdateOperation: &framework.PathOperation{
					Callback: handleSubaccountWrite,
					Summary:  "Create an UpCloud subaccount spec.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: handleSubaccountWrite,
					Summary:  "Create an UpCloud subaccount spec.",
				},
				// logical.DeleteOperation: &framework.PathOperation{
				// 	Callback: handleConfigDelete,
				// 	Summary:  "Deletes the secret at the specified location.",
				// },
				// logical.RevokeOperation: &framework.PathOperation{
				// 	Callback: ,
				// }
			},
			//ExistenceCheck: b.handleExistenceCheck,
		},
		{
			Pattern: "subaccount/" + framework.GenericNameRegex("subaccount") + "/token",

			Fields: map[string]*framework.FieldSchema{},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: handleSubaccountTokenRead,
					Summary:  "Retrieve the UpCloud subaccount credentials.",
				},
				// logical.DeleteOperation: &framework.PathOperation{
				// 	Callback: handleConfigDelete,
				// 	Summary:  "Deletes the secret at the specified location.",
				// },
				// logical.RevokeOperation: &framework.PathOperation{
				// 	Callback: ,
				// }
			},
			//ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func handleSubaccountWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	sa := data.Get("subaccount").(string)
	var subaccount upcloudSubaccount

	subaccount.Username = sa

	entry, err := logical.StorageEntryJSON("subaccount/"+sa, subaccount)
	if err != nil {
		return nil, fmt.Errorf("could not marshal json for upcloud auth: %w", err)
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, fmt.Errorf("could not put upcloudAuth to storage: %w", err)
	}

	return nil, nil
}

func handleSubaccountTokenRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// sa := data.Get("subaccount").(string)
	// var subaccount upcloudSubaccount

	// subaccount.Username = sa

	// entry, err := logical.StorageEntryJSON("subaccount/"+sa, subaccount)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not marshal json for upcloud auth: %w", err)
	// }

	// if err := req.Storage.Put(ctx, entry); err != nil {
	// 	return nil, fmt.Errorf("could not put upcloudAuth to storage: %w", err)
	// }

	return nil, nil
}
