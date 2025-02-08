package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Tsubasa-2005/GoMSGraph/graphhelper"
	"github.com/Tsubasa-2005/GoMSGraph/service"
)

func main() {
	clientId := os.Getenv("AZURE_CLIENT_ID")
	tenantId := os.Getenv("AZURE_TENANT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	driveID := os.Getenv("DRIVE_ID")
	parentFolderID := os.Getenv("DRIVE_ROOT_ITEM_ID")

	// siteName や filePath は .env に含めない場合、固定値や別途取得方法にしてください
	siteName := "YOUR_SITE_NAME"        // 適宜実際のサイト名に置き換えてください
	filePath := "path/to/your/file.txt" // アップロードするファイルのパス

	// オプション: カスタムロガーを利用する場合は作成、特に必要なければ nil でもOK
	var logger *graphhelper.Logger = nil

	// GraphService の初期化
	svc, err := service.NewGraphService(clientId, tenantId, clientSecret, logger)
	if err != nil {
		log.Fatalf("GraphService の初期化に失敗しました: %v", err)
	}

	ctx := context.Background()

	// 1. アプリケーションアクセストークンの取得例
	token, err := svc.GetAppToken(ctx)
	if err != nil {
		log.Printf("アクセストークンの取得エラー: %v", err)
	} else {
		fmt.Printf("取得したアクセストークン: %s\n", token.Token)
	}

	// 2. フォルダ作成の例
	folderName := "NewFolder"
	folderItem, err := svc.CreateFolder(ctx, driveID, parentFolderID, folderName)
	if err != nil {
		log.Printf("フォルダ作成エラー: %v", err)
	} else {
		fmt.Printf("作成したフォルダ: ID = %s, Name = %s\n", *folderItem.GetId(), *folderItem.GetName())
	}

	// 3. ドライブのルートアイテム一覧取得の例
	items, err := svc.GetDriveRootItems(ctx, driveID)
	if err != nil {
		log.Printf("ルートアイテム取得エラー: %v", err)
	} else {
		fmt.Println("ドライブのルートアイテム一覧:")
		for _, item := range items {
			fmt.Printf("ID: %s, Name: %s\n", *item.GetId(), *item.GetName())
		}
	}

	// 4. アイテム削除の例
	// ※ここでは先ほど作成したフォルダを削除する例を示します
	if folderItem != nil {
		err = svc.DeleteDriveItem(ctx, driveID, *folderItem.GetId())
		if err != nil {
			log.Printf("アイテム削除エラー: %v", err)
		} else {
			fmt.Println("フォルダの削除に成功しました")
		}
	}

	// 5. サイト検索の例
	sites, err := svc.GetSiteByName(ctx, siteName)
	if err != nil {
		log.Printf("サイト検索エラー: %v", err)
	} else {
		fmt.Printf("サイト '%s' の検索結果:\n", siteName)
		for _, s := range sites {
			// サイトのIDと表示名を出力
			fmt.Printf("ID: %s, Name: %s\n", *s.GetId(), *s.GetDisplayName())
		}
	}

	// 6. ファイルアップロードの例
	// 6-1. アップロードセッションの作成
	uploadSession, err := svc.CreateUploadSession(ctx, driveID, filePath)
	if err != nil {
		log.Printf("アップロードセッション作成エラー: %v", err)
	} else {
		// 6-2. アップロードするファイルをオープン
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("ファイルオープンエラー: %v", err)
		} else {
			defer file.Close()
			// 6-3. ファイルのアップロード
			uploadedItem, err := svc.UploadFile(uploadSession, file)
			if err != nil {
				log.Printf("ファイルアップロードエラー: %v", err)
			} else {
				fmt.Printf("アップロード完了: ID = %s, Name = %s\n", *uploadedItem.GetId(), *uploadedItem.GetName())
			}
		}
	}
}
