package graphhelper

import (
	"golang.org/x/net/context"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (g *GraphHelper) CreateFolder(ctx context.Context, driveID, driveItemID, folderName string) (models.DriveItemable, error) {
	newFolder := models.NewDriveItem()
	newFolder.SetName(&folderName)
	newFolder.SetFolder(models.NewFolder())

	return g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId(driveItemID).Children().Post(ctx, newFolder, nil)
}
