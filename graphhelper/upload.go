package graphhelper

import (
	"fmt"
	"os"
	"strings"

	"github.com/microsoftgraph/msgraph-sdk-go-core/fileuploader"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func (g *GraphHelper) UploadFile(uploadSession models.UploadSessionable, byteStream *os.File) (models.DriveItemable, error) {
	maxSliceSize := int64(320 * 1024)
	fileUploadTask := fileuploader.NewLargeFileUploadTask[models.DriveItemable](
		g.appClient.RequestAdapter,
		uploadSession,
		byteStream,
		maxSliceSize,
		models.CreateDriveItemFromDiscriminatorValue,
		nil,
	)

	progress := func(uploaded int64, total int64) {
		g.Logger.Debugf("Uploaded %d of %d bytes", uploaded, total)
	}

	uploadResult := fileUploadTask.Upload(progress)
	if uploadResult.GetUploadSucceeded() {
		g.Logger.Infof("Upload complete, item ID: %s", *uploadResult.GetItemResponse().GetId())
		return uploadResult.GetItemResponse(), nil
	}

	g.Logger.Warnf("Initial upload failed. Attempting to resume...")

	resumeResult, err := fileUploadTask.Resume(progress)
	if err != nil {
		g.Logger.Errorf("Upload resume failed: %v", err)
		return nil, fmt.Errorf("upload resume failed: %w", err)
	}
	if resumeResult.GetUploadSucceeded() {
		g.Logger.Infof("Upload resumed complete, item ID: %s", *resumeResult.GetItemResponse().GetId())
		return resumeResult.GetItemResponse(), nil
	}

	var errMessages []string
	for _, e := range resumeResult.GetResponseErrors() {
		errMessages = append(errMessages, e.Error())
	}
	joinedErrs := strings.Join(errMessages, "; ")
	g.Logger.Errorf("Upload failed: %s", joinedErrs)
	return nil, fmt.Errorf("upload failed: %s", joinedErrs)
}
