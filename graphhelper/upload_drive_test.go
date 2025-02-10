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
)

func TestGraphHelper_UploadFile(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)
	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID is not set as an environment variable, skipping test")
	}

	// Upload both a file size that will be chunked and a size that will not be chunked
	testCases := []struct {
		name     string
		fileSize int
	}{
		{name: "small file", fileSize: 100 * 1024},
		{name: "large file", fileSize: 400 * 1024},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpFile, err := os.CreateTemp("", fmt.Sprintf("upload_test_%s_*.dat", tc.name))
			assert.NoError(t, err, "Failed to create temporary file: %s", tc.name)
			defer os.Remove(tmpFile.Name())

			content := bytes.Repeat([]byte(tc.name[:1]), tc.fileSize)
			_, err = tmpFile.Write(content)
			assert.NoError(t, err, "Failed to write to file: %s", tc.name)
			tmpFile.Close()

			file, err := os.Open(tmpFile.Name())
			assert.NoError(t, err, "Failed to open file: %s", tc.name)
			defer file.Close()

			itemPath := fmt.Sprintf("test_upload/%s_%d.dat", tc.name, time.Now().UnixNano())
			uploadSession, err := gh.CreateUploadSession(context.Background(), driveID, itemPath)
			assert.NoError(t, err, "Failed to create an upload session: %s", tc.name)
			assert.NotNil(t, uploadSession.GetUploadUrl(), "UploadUrl is nil: %s", tc.name)

			driveItem, err := gh.UploadFile(uploadSession, file)
			assert.NoError(t, err, "Failed to execute UploadFile: %s", tc.name)
			assert.NotNil(t, driveItem, "The upload result is nil: %s", tc.name)
			if size := driveItem.GetSize(); size != nil {
				assert.Equal(t, int64(tc.fileSize), *size, "The size after upload does not match: %s", tc.name)
			}

			t.Logf("[%s] Uploaded drive item ID: %s", tc.name, *driveItem.GetId())
		})
	}
}
