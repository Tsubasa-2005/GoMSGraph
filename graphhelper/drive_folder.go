package graphhelper

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (g *GraphHelper) CreateFolder(ctx context.Context, driveID, driveItemID, folderName string) (models.DriveItemable, error) {
	newFolder := models.NewDriveItem()
	newFolder.SetName(&folderName)
	newFolder.SetFolder(models.NewFolder())

	item, err := g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Children().Post(ctx, newFolder, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder '%s' in drive '%s' at item '%s': %w", folderName, driveID, driveItemID, err)
	}

	return item, nil
}
