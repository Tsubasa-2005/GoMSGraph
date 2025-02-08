package graphhelper_test

import (
	"context"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGraphHelper_GetAppToken_Success(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	token, err := gh.GetAppToken(context.Background())
	if err != nil {
		t.Fatalf("GetAppToken の取得に失敗しました: %v", err)
	}

	assert.NotNil(t, token.Token, "取得したトークンが nil です")
	now := time.Now()
	assert.True(t, token.ExpiresOn.After(now) || token.ExpiresOn.Equal(now), "取得したトークンの有効期限が既に切れています。")
}
