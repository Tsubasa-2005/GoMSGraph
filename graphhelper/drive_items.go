package graphhelper

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (g *GraphHelper) GetDriveRootItems(ctx context.Context, driveID string) ([]models.DriveItemable, error) {
	get, err := g.appClient.Drives().ByDriveId(driveID).Root().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get drive root for drive %s: %w", driveID, err)
	}

	items, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(*get.GetId()).Children().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get children for drive item %s in drive %s: %w", *get.GetId(), driveID, err)
	}

	return items.GetValue(), nil
}

func (g *GraphHelper) GetDriveItem(ctx context.Context, driveID, driveItemID string) (models.DriveItemable, error) {
	item, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get drive item %s in drive %s: %w", driveItemID, driveID, err)
	}
	return item, nil
}

func (g *GraphHelper) DeleteDriveItem(ctx context.Context, driveID, driveItemID string) error {
	if err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Delete(ctx, nil); err != nil {
		return fmt.Errorf("failed to delete drive item %s in drive %s: %w", driveItemID, driveID, err)
	}
	return nil
}

func (g *GraphHelper) DownloadDriveItem(ctx context.Context, driveID, driveItemID string) ([]byte, error) {
	bytes, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Content().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download drive item %s in drive %s: %w", driveItemID, driveID, err)
	}
	return bytes, nil
}
