package graphhelper_test

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGraphHelper_CreateUploadSession(t *testing.T) {
	gh := testutil.SetUpGraphHelper(t)

	driveID := os.Getenv("DRIVE_ID")
	if driveID == "" {
		t.Skip("DRIVE_ID is not set in environment variables, skipping test")
	}

	res, err := gh.CreateUploadSession(context.Background(), driveID, "test/test.txt")
	if err != nil {
		t.Fatalf("Failed to call GetFileShareLink: %v", err)
	}

	assert.NotNil(t, res.GetUploadUrl())
	uploadUrl := res.GetUploadUrl()
	parsedUrl, err := url.Parse(*uploadUrl)
	assert.NoError(t, err, "Upload URL should be a valid URL")
	assert.NotNil(t, parsedUrl, "Parsed URL should not be nil")

	t.Logf("Upload URL: %s", *uploadUrl)
}
