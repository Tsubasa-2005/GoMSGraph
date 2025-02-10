package graphhelper

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func (g *GraphHelper) GetAppToken(ctx context.Context) (azcore.AccessToken, error) {
	token, err := g.clientSecretCredential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	})
	if err != nil {
		return azcore.AccessToken{}, fmt.Errorf("failed to get app token: %w", err)
	}

	return token, nil
}
