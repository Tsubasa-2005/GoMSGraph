package testutil

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/Tsubasa-2005/GoMSGraph/v2/graphhelper"
)

func SetUpGraphHelper(t *testing.T) *graphhelper.GraphHelper {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientID == "" || tenantID == "" || clientSecret == "" {
		t.Skip("Please set the environment variables AZURE_CLIENT_ID, AZURE_TENANT_ID, AZURE_CLIENT_SECRET")
	}
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to initialize Zap Logger: %v", err)
	}

	gh, err := graphhelper.NewGraphHelper(clientID, tenantID, clientSecret, graphhelper.NewDefaultLogger(zapLogger))
	if err != nil {
		t.Fatalf("Failed to initialize GraphHelper: %v", err)
	}
	return gh
}
