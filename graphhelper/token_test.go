package graphhelper_test

import (
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
)

func TestGraphHelper_GetAppToken_Success(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	token, err := gh.GetAppToken()
	if err != nil {
		t.Fatalf("GetAppToken の取得に失敗しました: %v", err)
	}

	if token.Token == "" {
		t.Errorf("取得したトークンが空です")
	}

	now := time.Now()
	if token.ExpiresOn.Before(now) {
		t.Errorf("取得したトークンの有効期限が既に切れています。Expires: %v, Now: %v", token.ExpiresOn, now)
	}

	t.Logf("取得したトークン: %s 有効期限: %v", token.Token, token.ExpiresOn)
}
