package graphhelper_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphHelper_CreateFolder(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID is not set as an environment variable, skipping test")
	}
	driveItemID := os.Getenv("DRIVE_ROOT_ITEM_ID")
	if driveItemID == "" {
		t.Skip("DRIVE_ITEM_ID is not set as an environment variable, skipping test")
	}

	folderName := "TestFolder_" + time.Now().Format("20060102150405")

	res, err := gh.CreateFolder(context.Background(), driveID, driveItemID, folderName)
	if err != nil {
		t.Fatalf("Failed to call CreateFolder: %v", err)
	}

	assert.NotNil(t, res.GetId(), "Folder ID should not be nil")
	assert.Equal(t, folderName, *res.GetName(), "Folder name should be equal")

	t.Cleanup(func() {
		t.Cleanup(func() {
			require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, *res.GetId()), "作成したフォルダの削除に失敗しました")
		})
	})
}
