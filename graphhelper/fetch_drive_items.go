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
