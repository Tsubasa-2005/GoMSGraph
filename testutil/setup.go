package testutil

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/Tsubasa-2005/GoMSGraph/graphhelper"
	"github.com/Tsubasa-2005/GoMSGraph/httpclient"
)

func SetUpGraphHelper(t *testing.T) *graphhelper.GraphHelper {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientID == "" || tenantID == "" || clientSecret == "" {
		t.Skip("環境変数 AZURE_CLIENT_ID, AZURE_TENANT_ID, AZURE_CLIENT_SECRET を設定してください")
	}
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Zap Logger の初期化に失敗しました: %v", err)
	}

	gh, err := graphhelper.NewGraphHelper(clientID, tenantID, clientSecret, graphhelper.NewDefaultLogger(zapLogger))
	if err != nil {
		t.Fatalf("GraphHelper の初期化に失敗しました: %v", err)
	}
	return gh
}

func SetUpHttpClient(t *testing.T) *httpclient.Client {
	baseURL := os.Getenv("GRAPH_API_BASE_URL")
	gh := SetUpGraphHelper(t)
	token, err := gh.GetAppToken()
	if err != nil {
		t.Fatalf("アクセストークンの取得に失敗しました: %v", err)
	}

	return httpclient.NewClient(baseURL, token.Token, nil)
}
