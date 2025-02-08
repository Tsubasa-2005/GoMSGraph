package graphhelper_test

import (
	"context"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGraphHelper_GetAppToken_Success(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	token, err := gh.GetAppToken(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve GetAppToken: %v", err)
	}

	assert.NotNil(t, token.Token, "The retrieved token is nil")
	now := time.Now()
	assert.True(t, token.ExpiresOn.After(now) || token.ExpiresOn.Equal(now), "The retrieved token is already expired.")
}
