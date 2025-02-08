package graphhelper_test

import (
	"context"
	"os"
	"testing"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
)

func TestGraphHelper_GetDriveRootItems(t *testing.T) {
	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID が環境変数に設定されていないため、テストをスキップします")
	}

	_, err := gh.GetDriveRootItems(context.Background(), driveID)
	if err != nil {
		t.Fatalf("GetDriveRootItems の呼び出しに失敗しました: %v", err)
	}
}
