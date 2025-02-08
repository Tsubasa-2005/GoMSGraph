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
	gh := testutil.SetUpGraphHelper(t)
	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID が環境変数に設定されていないため、テストをスキップします")
	}

	// 分割されるサイズとそうでないサイズのファイルをアップロード
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
			assert.NoError(t, err, "一時ファイルの作成に失敗しました: %s", tc.name)
			defer os.Remove(tmpFile.Name())

			content := bytes.Repeat([]byte(tc.name[:1]), tc.fileSize)
			_, err = tmpFile.Write(content)
			assert.NoError(t, err, "ファイルへの書き込みに失敗しました: %s", tc.name)
			tmpFile.Close()

			file, err := os.Open(tmpFile.Name())
			assert.NoError(t, err, "ファイルのオープンに失敗しました: %s", tc.name)
			defer file.Close()

			itemPath := fmt.Sprintf("test_upload/%s_%d.dat", tc.name, time.Now().UnixNano())
			uploadSession, err := gh.CreateUploadSession(context.Background(), driveID, itemPath)
			assert.NoError(t, err, "アップロードセッションの作成に失敗しました: %s", tc.name)
			assert.NotNil(t, uploadSession.GetUploadUrl(), "UploadUrl が nil です: %s", tc.name)

			driveItem, err := gh.UploadFile(uploadSession, file)
			assert.NoError(t, err, "UploadFile の実行に失敗しました: %s", tc.name)
			assert.NotNil(t, driveItem, "アップロード結果が nil です: %s", tc.name)
			if size := driveItem.GetSize(); size != nil {
				assert.Equal(t, int64(tc.fileSize), *size, "アップロード後のサイズが一致しません: %s", tc.name)
			}

			t.Logf("[%s] アップロードされたドライブアイテムの ID: %s", tc.name, *driveItem.GetId())
		})
	}
}
