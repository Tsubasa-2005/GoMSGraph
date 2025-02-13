package graphhelper_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/v2/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphHelper_GetDriveRootItems(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID is not set in environment variables, skipping test")
	}

	_, err := gh.GetDriveRootItems(context.Background(), driveID)
	if err != nil {
		t.Fatalf("Failed to call GetDriveRootItems: %v", err)
	}
}

func TestGraphHelper_GetDriveItem_Multiple(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)
	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_IDが環境変数に設定されていないため、テストをスキップします")
	}
	rootItemID := os.Getenv("DRIVE_ROOT_ITEM_ID")
	if rootItemID == "" {
		t.Skip("DRIVE_ROOT_ITEM_IDが環境変数に設定されていないため、テストをスキップします")
	}

	tests := []struct {
		name        string
		itemType    string
		setup       func(t *testing.T) string
		expectError bool
	}{
		{
			name:     "File",
			itemType: "file",
			setup: func(t *testing.T) string {
				tmpFile, err := os.CreateTemp("", "test_file_*.dat")
				require.NoError(t, err, "一時ファイルの作成に失敗しました")
				defer os.Remove(tmpFile.Name())

				content := []byte("This is test file content")
				_, err = tmpFile.Write(content)
				require.NoError(t, err, "一時ファイルへの書き込みに失敗しました")
				tmpFile.Close()

				file, err := os.Open(tmpFile.Name())
				require.NoError(t, err, "一時ファイルのオープンに失敗しました")
				defer file.Close()

				fileName := fmt.Sprintf("TestFile_%d.dat", time.Now().UnixNano())
				uploadSession, err := gh.CreateUploadSession(context.Background(), driveID, fileName)
				require.NoError(t, err, "アップロードセッションの作成に失敗しました")
				require.NotNil(t, uploadSession.GetUploadUrl(), "UploadUrlがnilです")

				driveItem, err := gh.UploadFile(uploadSession, file)
				require.NoError(t, err, "ファイルのアップロードに失敗しました")
				require.NotNil(t, driveItem, "アップロード結果がnilです")
				require.NotNil(t, driveItem.GetId(), "アップロードしたファイルのIDがnilです")

				t.Cleanup(func() {
					require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, *driveItem.GetId()), "アップロードしたファイルの削除に失敗しました")
				})

				return *driveItem.GetId()
			},
			expectError: false,
		},
		{
			name:     "Folder",
			itemType: "folder",
			setup: func(t *testing.T) string {
				folderName := "TestFolder_" + time.Now().Format("20060102150405")
				folderItem, err := gh.CreateFolder(context.Background(), driveID, rootItemID, folderName)
				require.NoError(t, err, "フォルダの作成に失敗しました")
				require.NotNil(t, folderItem, "作成したフォルダがnilです")
				require.NotNil(t, folderItem.GetId(), "作成したフォルダのIDがnilです")

				t.Cleanup(func() {
					require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, *folderItem.GetId()), "作成したフォルダの削除に失敗しました")
				})

				return *folderItem.GetId()
			},
			expectError: false,
		},
		{
			name:     "InvalidID",
			itemType: "invalid",
			setup: func(t *testing.T) string {
				return "non-existent-id"
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			driveItemID := tc.setup(t)

			item, err := gh.GetDriveItem(context.Background(), driveID, driveItemID)
			if tc.expectError {
				require.Error(t, err, "エラーが発生することを期待します: %s", tc.name)
				assert.Nil(t, item, "取得結果はnilであるべきです: %s", tc.name)
			} else {
				require.NoError(t, err, "エラーが発生してはいけません: %s", tc.name)
				require.NotNil(t, item, "取得したDriveItemがnilです: %s", tc.name)
				assert.Equal(t, driveItemID, *item.GetId(), "取得したDriveItemのIDが一致しません: %s", tc.name)
			}
		})
	}
}

func TestGraphHelper_DeleteDriveItem(t *testing.T) {
	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID is not set in environment variables, skipping test")
	}
	rootItemID := os.Getenv("DRIVE_ROOT_ITEM_ID")
	if rootItemID == "" {
		t.Skip("DRIVE_ROOT_ITEM_ID is not set in environment variables, skipping test")
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
				require.NoError(t, err, "Failed to create folder")
				require.NotNil(t, createRes.GetId(), "The created folder ID is nil")
				folderID := *createRes.GetId()

				require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, folderID), "Failed to delete folder")

				items, err := gh.GetDriveRootItems(context.Background(), driveID)
				require.NoError(t, err, "Failed to retrieve root item list")
				found := false
				for _, item := range items {
					if item.GetId() != nil && *item.GetId() == folderID {
						found = true
						break
					}
				}
				require.False(t, found, "The deleted folder still exists in the root item list")

			case "file":
				tmpFile, err := os.CreateTemp("", "upload_test_file_*.dat")
				assert.NoError(t, err, "Failed to create a temporary file")
				defer os.Remove(tmpFile.Name())

				fileSize := 100 * 1024
				content := bytes.Repeat([]byte("A"), fileSize)
				_, err = tmpFile.Write(content)
				assert.NoError(t, err, "Failed to write to temporary file")
				tmpFile.Close()

				file, err := os.Open(tmpFile.Name())
				assert.NoError(t, err, "Failed to open temporary file")
				defer file.Close()

				fileName := fmt.Sprintf("TestFileToDelete_%d.dat", time.Now().UnixNano())
				uploadSession, err := gh.CreateUploadSession(context.Background(), driveID, fileName)
				assert.NoError(t, err, "Failed to create upload session")
				assert.NotNil(t, uploadSession.GetUploadUrl(), "UploadUrl is nil")

				driveItem, err := gh.UploadFile(uploadSession, file)
				assert.NoError(t, err, "Failed to upload file")
				assert.NotNil(t, driveItem, "The upload result is nil")
				fileID := *driveItem.GetId()

				err = gh.DeleteDriveItem(context.Background(), driveID, fileID)
				assert.NoError(t, err, "Failed to delete file")

				items, err := gh.GetDriveRootItems(context.Background(), driveID)
				assert.NoError(t, err, "Failed to retrieve root item list")
				found := false
				for _, item := range items {
					if item.GetId() != nil && *item.GetId() == fileID {
						found = true
						break
					}
				}
				assert.False(t, found, "The deleted file still exists in the root item list")
			}
		})
	}
}

func TestGraphHelper_DownloadDriveItem(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_IDが環境変数に設定されていないため、テストをスキップします")
	}
	rootItemID := os.Getenv("DRIVE_ROOT_ITEM_ID")
	if rootItemID == "" {
		t.Skip("DRIVE_ROOT_ITEM_IDが環境変数に設定されていないため、テストをスキップします")
	}

	tests := []struct {
		name     string
		itemType string
	}{
		{name: "File Download", itemType: "file"},
		{name: "Folder Download", itemType: "folder"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			switch tc.itemType {
			case "file":
				tmpFile, err := os.CreateTemp("", "download_test_file_*.dat")
				require.NoError(t, err, "一時ファイルの作成に失敗しました")
				defer os.Remove(tmpFile.Name())

				fileSize := 100 * 1024
				content := bytes.Repeat([]byte("B"), fileSize)
				_, err = tmpFile.Write(content)
				require.NoError(t, err, "一時ファイルへの書き込みに失敗しました")
				tmpFile.Close()

				file, err := os.Open(tmpFile.Name())
				require.NoError(t, err, "一時ファイルのオープンに失敗しました")
				defer file.Close()

				fileName := fmt.Sprintf("TestFileDownload_%d.dat", time.Now().UnixNano())
				uploadSession, err := gh.CreateUploadSession(context.Background(), driveID, fileName)
				require.NoError(t, err, "アップロードセッションの作成に失敗しました")
				require.NotNil(t, uploadSession.GetUploadUrl(), "UploadUrlがnilです")

				driveItem, err := gh.UploadFile(uploadSession, file)
				require.NoError(t, err, "ファイルのアップロードに失敗しました")
				require.NotNil(t, driveItem, "アップロード結果がnilです")
				fileID := *driveItem.GetId()

				t.Cleanup(func() {
					require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, fileID), "アップロードしたファイルの削除に失敗しました")
				})

				downloadedContent, err := gh.DownloadDriveItem(context.Background(), driveID, fileID)
				require.NoError(t, err, "ファイルのダウンロードに失敗しました")
				assert.Equal(t, content, downloadedContent, "ダウンロードした内容が元の内容と一致しません")

			case "folder":
				folderName := "TestFolderDownload_" + time.Now().Format("20060102150405")
				folderItem, err := gh.CreateFolder(context.Background(), driveID, rootItemID, folderName)
				require.NoError(t, err, "フォルダの作成に失敗しました")
				require.NotNil(t, folderItem.GetId(), "作成したフォルダのIDがnilです")
				folderID := *folderItem.GetId()

				t.Cleanup(func() {
					require.NoError(t, gh.DeleteDriveItem(context.Background(), driveID, folderID), "作成したフォルダの削除に失敗しました")
				})

				downloadedContent, err := gh.DownloadDriveItem(context.Background(), driveID, folderID)
				require.Error(t, err, "フォルダのダウンロード時にエラーが発生することを期待しています")
				assert.Nil(t, downloadedContent, "フォルダのダウンロード結果はnilであるべきです")
			}
		})
	}
}
