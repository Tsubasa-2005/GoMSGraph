package graphhelper

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphHelper struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
	Logger                 *Logger
}

func NewGraphHelper(clientId, tenantId, clientSecret string, logger *Logger) (*GraphHelper, error) {
	credential, authProvider, err := initializeAuth(clientId, tenantId, clientSecret)
	if err != nil {
		return nil, err
	}
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create request adapter: %w", err)
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)

	if logger == nil {
		zapLogger, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create default zap logger: %w", err)
		}
		logger = NewDefaultLogger(zapLogger)
	}

	return &GraphHelper{
		clientSecretCredential: credential,
		appClient:              client,
		Logger:                 logger,
	}, nil
}

func initializeAuth(clientId, tenantId, clientSecret string) (*azidentity.ClientSecretCredential, *auth.AzureIdentityAuthenticationProvider, error) {
	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create an authentication provider: %w", err)
	}

	return credential, authProvider, nil
}
