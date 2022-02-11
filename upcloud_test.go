package upcloud

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/service"
	"github.com/stretchr/testify/require"
)

func TestUpcloudAPI(t *testing.T) {
	username := "jlarfors"
	password := "<REDACTED>"

	svc := service.New(client.New(username, password))

	_, err := svc.CreateSubaccount(&request.CreateSubaccountRequest{
		Subaccount: request.CreateSubaccount{
			Email:    "noreply@verifa.io",
			Phone:    "+358.400000000",
			Username: "hacking",
			Password: "SomeSuperSecret123$$",
			AllowAPI: 1,
			Language: "en",
			Currency: "EUR",
			Timezone: "UTC",
			// FirstName: "Mr",
			// LastName: "Hack",
			// Company: "verifa",
			// Email: ,
		},
	})
	require.NoError(t, err)
}
