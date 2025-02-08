package graphhelper

import (
	"context"

	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// CreateUploadSession は、指定されたドライブにアイテムをアップロードするためのセッションを作成します。
// itemPathで指定されたパスにアイテムが存在しない場合は、新規にアイテムが作成されます。
func (g *GraphHelper) CreateUploadSession(ctx context.Context, driveID, itemPath string) (models.UploadSessionable, error) {
	itemUploadProperties := models.NewDriveItemUploadableProperties()
	itemUploadProperties.SetAdditionalData(map[string]any{"@microsoft.graph.conflictBehavior": "replace"})
	uploadSessionRequestBody := drives.NewItemItemsItemCreateUploadSessionPostRequestBody()
	uploadSessionRequestBody.SetItem(itemUploadProperties)

	return g.appClient.Drives().ByDriveId(driveID).Items().ByDriveItemId("root:/"+itemPath+":").CreateUploadSession().Post(ctx, uploadSessionRequestBody, nil)
}
