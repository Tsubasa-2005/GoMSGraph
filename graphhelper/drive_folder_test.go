package graphhelper_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGraphHelper_CreateFolder(t *testing.T) {
	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID が環境変数に設定されていないため、テストをスキップします")
	}
	driveItemID := os.Getenv("DRIVE_ROOT_ITEM_ID")
	if driveItemID == "" {
		t.Skip("DRIVE_ITEM_ID が環境変数に設定されていないため、テストをスキップします")
	}

	folderName := "TestFolder_" + time.Now().Format("20060102150405")

	res, err := gh.CreateFolder(context.Background(), driveID, driveItemID, folderName)
	if err != nil {
		t.Fatalf("CreateFolder の呼び出しに失敗しました: %v", err)
	}

	assert.NotNil(t, res.GetId(), "Folder ID should not be nil")
	assert.Equal(t, folderName, *res.GetName(), "Folder name should be equal")
}
