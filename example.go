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

	// If you do not include siteName or filePath in .env, set fixed values or implement another method to retrieve them
	siteName := "YOUR_SITE_NAME"        // Replace with the actual site name as needed
	filePath := "path/to/your/file.txt" // Path to the file to be uploaded

	// Optional: If you want to use a custom logger, create it here. If not necessary, you can set it to nil.
	var logger *graphhelper.Logger = nil

	// Initialize GraphService
	svc, err := service.NewGraphService(clientId, tenantId, clientSecret, logger)
	if err != nil {
		log.Fatalf("Failed to initialize GraphService: %v", err)
	}

	ctx := context.Background()

	// 1. Example of retrieving an application access token
	token, err := svc.GetAppToken(ctx)
	if err != nil {
		log.Printf("Error retrieving access token: %v", err)
	} else {
		fmt.Printf("Retrieved access token: %s\n", token.Token)
	}

	// 2. Example of creating a folder
	folderName := "NewFolder"
	folderItem, err := svc.CreateFolder(ctx, driveID, parentFolderID, folderName)
	if err != nil {
		log.Printf("Folder creation error: %v", err)
	} else {
		fmt.Printf("Folder created: ID = %s, Name = %s\n", *folderItem.GetId(), *folderItem.GetName())
	}

	// 3. Example of retrieving drive root items
	items, err := svc.GetDriveRootItems(ctx, driveID)
	if err != nil {
		log.Printf("Root item retrieval error: %v", err)
	} else {
		fmt.Println("List of drive root items:")
		for _, item := range items {
			fmt.Printf("ID: %s, Name: %s\n", *item.GetId(), *item.GetName())
		}
	}

	// 4. Example of deleting an item
	// ※ This shows an example of deleting the folder you just created
	if folderItem != nil {
		err = svc.DeleteDriveItem(ctx, driveID, *folderItem.GetId())
		if err != nil {
			log.Printf("Item deletion error: %v", err)
		} else {
			fmt.Println("Successfully deleted the folder")
		}
	}

	// 5. Example of searching for a site
	sites, err := svc.GetSiteByName(ctx, siteName)
	if err != nil {
		log.Printf("Site search error: %v", err)
	} else {
		fmt.Printf("Search results for site '%s':\n", siteName)
		for _, s := range sites {
			// Output the site's ID and display name
			fmt.Printf("ID: %s, Name: %s\n", *s.GetId(), *s.GetDisplayName())
		}
	}

	// 6. Example of uploading a file
	// 6-1. Create an upload session
	uploadSession, err := svc.CreateUploadSession(ctx, driveID, filePath)
	if err != nil {
		log.Printf("Upload session creation error: %v", err)
	} else {
		// 6-2. Open the file to be uploaded
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("File open error: %v", err)
		} else {
			defer file.Close()
			// 6-3. Upload the file
			uploadedItem, err := svc.UploadFile(uploadSession, file)
			if err != nil {
				log.Printf("File upload error: %v", err)
			} else {
				fmt.Printf("Upload complete: ID = %s, Name = %s\n", *uploadedItem.GetId(), *uploadedItem.GetName())
			}
		}
	}
}
