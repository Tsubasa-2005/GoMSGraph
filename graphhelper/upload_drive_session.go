package graphhelper

import (
	"context"

	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// CreateUploadSession creates a session to upload an item to the specified drive.
// If no item exists at the path specified by itemPath, a new item will be created.
func (g *GraphHelper) CreateUploadSession(ctx context.Context, driveID, itemPath string) (models.UploadSessionable, error) {
	itemUploadProperties := models.NewDriveItemUploadableProperties()
	itemUploadProperties.SetAdditionalData(map[string]any{"@microsoft.graph.conflictBehavior": "replace"})
	uploadSessionRequestBody := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	uploadSessionRequestBody.SetItem(itemUploadProperties)

	return g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId("root:/"+itemPath+":").CreateUploadSession().Post(ctx, uploadSessionRequestBody, nil)
}
