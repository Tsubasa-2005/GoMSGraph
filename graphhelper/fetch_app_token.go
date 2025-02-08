package graphhelper

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func (g *GraphHelper) GetAppToken(ctx context.Context) (azcore.AccessToken, error) {
	return g.clientSecretCredential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
}
