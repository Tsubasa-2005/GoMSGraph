package graphhelper_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGraphHelper_DeleteDriveItem(t *testing.T) {
	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID が環境変数に設定されていないため、テストをスキップします")
	}
	rootItemID := os.Getenv("DRIVE_ROOT_ITEM_ID")
	if rootItemID == "" {
		t.Skip("DRIVE_ROOT_ITEM_ID が環境変数に設定されていないため、テストをスキップします")
	}

	testCases := []struct {
		name     string
		itemType string
	}{
		{name: "Folder Deletion", itemType: "folder"},
		{name: "File Deletion", itemType: "file"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			switch tc.itemType {
			case "folder":
				folderName := "TestFolderToDelete_" + time.Now().Format("20060102150405")
				createRes, err := gh.CreateFolder(context.Background(), driveID, rootItemID, folderName)
				require.NoError(t, err, "フォルダ作成に失敗しました")
				require.NotNil(t, createRes.GetId(), "作成されたフォルダのIDが nil です")
				folderID := *createRes.GetId()

				require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, folderID), "フォルダ削除に失敗しました")

				items, err := gh.GetDriveRootItems(context.Background(), driveID)
				require.NoError(t, err, "ルートアイテム一覧の取得に失敗しました")
				found := false
				for _, item := range items {
					if item.GetId() != nil && *item.GetId() == folderID {
						found = true
						break
					}
				}
				require.False(t, found, "削除されたフォルダがルートアイテム一覧に存在しています")

			case "file":
				tmpFile, err := os.CreateTemp("", "upload_test_file_*.dat")
				assert.NoError(t, err, "一時ファイルの作成に失敗しました")
				defer os.Remove(tmpFile.Name())

				fileSize := 100 * 1024
				content := bytes.Repeat([]byte("A"), fileSize)
				_, err = tmpFile.Write(content)
				assert.NoError(t, err, "一時ファイルへの書き込みに失敗しました")
				tmpFile.Close()

				file, err := os.Open(tmpFile.Name())
				assert.NoError(t, err, "一時ファイルのオープンに失敗しました")
				defer file.Close()

				fileName := fmt.Sprintf("TestFileToDelete_%d.dat", time.Now().UnixNano())
				uploadSession, err := gh.CreateUploadSession(context.Background(), driveID, fileName)
				assert.NoError(t, err, "アップロードセッションの作成に失敗しました")
				assert.NotNil(t, uploadSession.GetUploadUrl(), "UploadUrl が nil です")

				driveItem, err := gh.UploadFile(uploadSession, file)
				assert.NoError(t, err, "ファイルのアップロードに失敗しました")
				assert.NotNil(t, driveItem, "アップロード結果が nil です")
				fileID := *driveItem.GetId()

				err = gh.DeleteDriveItem(context.Background(), driveID, fileID)
				assert.NoError(t, err, "ファイル削除に失敗しました")

				items, err := gh.GetDriveRootItems(context.Background(), driveID)
				assert.NoError(t, err, "ルートアイテム一覧の取得に失敗しました")
				found := false
				for _, item := range items {
					if item.GetId() != nil && *item.GetId() == fileID {
						found = true
						break
					}
				}
				assert.False(t, found, "削除されたファイルがルートアイテム一覧に存在しています")
			}
		})
	}
}
