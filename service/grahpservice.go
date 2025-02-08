package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Tsubasa-2005/GoMSGraph/graphhelper"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type GraphService interface {
	CreateFolder(ctx context.Context, driveID, driveItemID, folderName string) (models.DriveItemable, error)
	GetDriveRootItems(ctx context.Context, driveID string) ([]models.DriveItemable, error)
	DeleteDriveItem(ctx context.Context, driveID, driveItemID string) error
	GetAppToken(ctx context.Context) (azcore.AccessToken, error)
	GetSiteByName(ctx context.Context, siteName string) ([]models.Siteable, error)
	UploadFile(uploadSession models.UploadSessionable, file *os.File) (models.DriveItemable, error)
	CreateUploadSession(ctx context.Context, driveID, itemPath string) (models.UploadSessionable, error)
}

type graphServiceImpl struct {
	helper *graphhelper.GraphHelper
}

func NewGraphService(clientId, tenantId, clientSecret string, logger *graphhelper.Logger) (GraphService, error) {
	gh, err := graphhelper.NewGraphHelper(clientId, tenantId, clientSecret, logger)
	if err != nil {
		return nil, fmt.Errorf("GraphHelper の初期化に失敗しました: %w", err)
	}
	return &graphServiceImpl{
		helper: gh,
	}, nil
}

func (s *graphServiceImpl) CreateFolder(ctx context.Context, driveID, driveItemID, folderName string) (models.DriveItemable, error) {
	return s.helper.CreateFolder(ctx, driveID, driveItemID, folderName)
}

func (s *graphServiceImpl) GetDriveRootItems(ctx context.Context, driveID string) ([]models.DriveItemable, error) {
	return s.helper.GetDriveRootItems(ctx, driveID)
}

func (s *graphServiceImpl) DeleteDriveItem(ctx context.Context, driveID, driveItemID string) error {
	return s.helper.DeleteDriveItem(ctx, driveID, driveItemID)
}

func (s *graphServiceImpl) GetAppToken(ctx context.Context) (azcore.AccessToken, error) {
	return s.helper.GetAppToken(ctx)
}

func (s *graphServiceImpl) GetSiteByName(ctx context.Context, siteName string) ([]models.Siteable, error) {
	return s.helper.GetSiteByName(ctx, siteName)
}

func (s *graphServiceImpl) UploadFile(uploadSession models.UploadSessionable, file *os.File) (models.DriveItemable, error) {
	return s.helper.UploadFile(uploadSession, file)
}

func (s *graphServiceImpl) CreateUploadSession(ctx context.Context, driveID, itemPath string) (models.UploadSessionable, error) {
	return s.helper.CreateUploadSession(ctx, driveID, itemPath)
}
