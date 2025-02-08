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
		t.Skip("DRIVE_ID is not set in environment variables, skipping test")
	}

	_, err := gh.GetDriveRootItems(context.Background(), driveID)
	if err != nil {
		t.Fatalf("Failed to call GetDriveRootItems: %v", err)
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
