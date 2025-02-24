package graphhelper_test

import (
	"context"
	"testing"
	"time"

	"github.com/Tsubasa-2005/GoMSGraph/v2/testutil"
	"github.com/stretchr/testify/require"
)

func TestGraphHelper_GetSiteByName(t *testing.T) {
	t.Parallel()

	gh := testutil.SetUpGraphHelper(t)

	tests := []struct {
		name        string
		siteName    string
		expectError bool
	}{
		{
			name:        "SuccessCase",
			siteName:    "ITO_TestTeam_202411",
			expectError: false,
		},
		{
			name:        "NotFoundCase",
			siteName:    "ThisSiteDoesNotExist_12345",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			sites, err := gh.GetSiteByName(ctx, tt.siteName)
			if tt.expectError {
				require.Error(t, err, "Expected an error but did not get one")
			} else {
				if err != nil {
					t.Fatalf("Failed to retrieve site name '%s': %v", tt.siteName, err)
				}
				if len(sites) == 0 {
					t.Fatalf("No site found with site name '%s'", tt.siteName)
				}
				for _, site := range sites {
					id := site.GetId()
					displayName := site.GetDisplayName()
					var idStr, nameStr string
					if id != nil {
						idStr = *id
					}
					if displayName != nil {
						nameStr = *displayName
					}
					t.Logf("Site ID: %s, Display Name: %s", idStr, nameStr)
				}
			}
		})
	}
}
