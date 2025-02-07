package graphhelper_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
)

func TestGraphHelper_GetSiteByName(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	tests := []struct {
		name        string
		siteName    string
		expectError bool
	}{
		{
			name:        "SuccessCase",
			siteName:    "ITO_TestTeam_202411",
			expectError: false,
		},
		{
			name:        "NotFoundCase",
			siteName:    "ThisSiteDoesNotExist_12345",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			sites, err := gh.GetSiteByName(ctx, tt.siteName)
			if tt.expectError {
				expectedErrMsg := fmt.Sprintf("検索クエリ '%s' に一致するサイトが見つかりませんでした", tt.siteName)
				if err.Error() != expectedErrMsg {
					t.Errorf("エラーメッセージが期待値と異なります。期待: %q, 実際: %q", expectedErrMsg, err.Error())
				} else {
					t.Logf("期待通りエラーメッセージが返されました: %q", err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("サイト名 '%s' の取得に失敗しました: %v", tt.siteName, err)
				}
				if len(sites) == 0 {
					t.Fatalf("サイト名 '%s' でサイトが見つかりませんでした", tt.siteName)
				}
				for _, site := range sites {
					id := site.GetId()
					displayName := site.GetDisplayName()
					var idStr, nameStr string
					if id != nil {
						idStr = *id
					}
					if displayName != nil {
						nameStr = *displayName
					}
					t.Logf("Site ID: %s, Display Name: %s", idStr, nameStr)
				}
			}
		})
	}
}
