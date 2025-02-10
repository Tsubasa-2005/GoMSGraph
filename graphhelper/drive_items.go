package graphhelper

import (
	"context"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (g *GraphHelper) GetDriveRootItems(ctx context.Context, driveID string) ([]models.DriveItemable, error) {
	get, err := g.appClient.Drives().ByDriveId(driveID).Root().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	items, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(*get.GetId()).Children().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return items.GetValue(), nil
}

func (g *GraphHelper) GetDriveItem(ctx context.Context, driveID, driveItemID string) (models.DriveItemable, error) {
	item, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (g *GraphHelper) DeleteDriveItem(ctx context.Context, driveID, driveItemID string) error {
	return g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Delete(ctx, nil)
}

func (g *GraphHelper) DownloadDriveItem(ctx context.Context, driveID, driveItemID string) ([]byte, error) {
	bytes, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Content().Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
