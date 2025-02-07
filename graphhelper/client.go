package graphhelper

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphHelper struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
}

func NewGraphHelper(clientId, tenantId, clientSecret string) (*GraphHelper, error) {
	credential, authProvider, err := initializeAuth(clientId, tenantId, clientSecret)
	if err != nil {
		return nil, err
	}
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return nil, fmt.Errorf("リクエストアダプターの作成に失敗しました: %w", err)
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)
	return &GraphHelper{
		clientSecretCredential: credential,
		appClient:              client,
	}, nil
}

func initializeAuth(clientId, tenantId, clientSecret string) (*azidentity.ClientSecretCredential, *auth.AzureIdentityAuthenticationProvider, error) {
	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("認証情報の作成に失敗しました: %w", err)
	}

	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("認証プロバイダーの作成に失敗しました: %w", err)
	}

	return credential, authProvider, nil
}
