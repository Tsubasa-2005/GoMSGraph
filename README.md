# GoMSGraph

GoMSGraph は、Microsoft Graph API を簡単に利用するための Go 言語ライブラリです。  
Azure AD アプリケーションを用いた認証処理から、ドライブアイテム(フォルダ/ファイル)の操作、SharePoint サイトの検索、大容量ファイルアップロードなどをシンプルなメソッドで実行できます。

---

## 特徴

- **認証処理の簡略化**  
  Azure SDK の `azidentity` を利用し、Microsoft Graph API への認証をシンプルに実装。  
  `GetAppToken` メソッドを通じてアクセストークンを容易に取得可能です。

- **ドライブ操作**
  - **フォルダの作成**: `CreateFolder` メソッドで指定のドライブにフォルダを作成
  - **ルートアイテム一覧の取得**: `GetDriveRootItems` メソッドでドライブのルートアイテムを取得
  - **アイテム削除**: `DeleteDriveItem` メソッドで任意のアイテム (ファイル/フォルダ) を削除

- **サイト検索**  
  指定したサイト名で、Microsoft 365 内のサイトを検索する `GetSiteByName` メソッドを提供。

- **大容量ファイルアップロード**  
  アップロードセッションを作成する `CreateUploadSession` と、  
  セッションを使った分割アップロードの `UploadFile` メソッドで、大容量ファイルの取り扱いに対応。

- **カスタムロガー対応**  
  デフォルトロガーには[Uber Zap](https://github.com/uber-go/zap)を用いており、必要に応じてカスタムロガーの差し替えが可能です。

---

## ディレクトリ構成

ライブラリ内の主なファイル構成は以下の通りです。

```
graphhelper
├── client.go
├── drive_folder.go
├── drive_folder_test.go
├── drive_items.go
├── drive_items_test.go
├── fetch_app_token.go
├── fetch_app_token_test.go
├── fetch_sites.go
├── fetch_sites_test.go
├── upload_drive.go
├── upload_drive_session.go
├── upload_drive_session_test.go
├── upload_drive_test.go
└── zaplogger.go
```

```
service
└── grahpservice.go
```

---

## 必要条件

- Go 1.18 以上
- Microsoft Graph API へのアクセス権を持つ Azure アプリケーション (以下を取得済みであること)
  - クライアント ID (`AZURE_CLIENT_ID`)
  - テナント ID (`AZURE_TENANT_ID`)
  - クライアントシークレット (`AZURE_CLIENT_SECRET`)
- Microsoft Graph API で使用するスコープ： `https://graph.microsoft.com/.default`
- 操作対象の OneDrive または SharePoint ドライブ ID (`DRIVE_ID`)  
  ※ ルートフォルダ以外を指定したい場合は、そのアイテムID (`DRIVE_ROOT_ITEM_ID`) も必要です。

---

## インストール

Go Modules を使っているプロジェクトにて、以下のコマンドを実行してください。

```bash
go get github.com/Tsubasa-2005/GoMSGraph
```

---

## 環境変数の設定

このライブラリを使用するには、認証情報やドライブ情報を環境変数で設定します。

- `AZURE_CLIENT_ID`  
  Azure AD アプリのクライアント ID
- `AZURE_TENANT_ID`  
  Azure テナント ID
- `AZURE_CLIENT_SECRET`  
  Azure AD アプリのクライアントシークレット
- `DRIVE_ID`  
  操作対象となるドライブ ID
- `DRIVE_ROOT_ITEM_ID`  
  ルートフォルダや任意フォルダのアイテム ID (フォルダ作成時の親 ID などに利用)

例:
```bash
export AZURE_CLIENT_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
export AZURE_TENANT_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
export AZURE_CLIENT_SECRET="xxxxxxxxxxxxxxxxxxxx"
export DRIVE_ID="xxxxxxxxxxxxxxxxxxxx"
export DRIVE_ROOT_ITEM_ID="xxxxxxxxxxxxxxxxxxxx"
```

---

## 使い方

以下は、本ライブラリの主な機能を利用するためのサンプルコード例です。  
[example.go](./example.go)


---

## 主なメソッド一覧

- **認証 / トークン取得**
  - `GetAppToken(ctx context.Context) (azcore.AccessToken, error)`  
    アプリケーション権限で Microsoft Graph へアクセスするためのアクセストークンを取得します。

- **ドライブ操作**
  - `CreateFolder(ctx context.Context, driveID, driveItemID, folderName string) (models.DriveItemable, error)`  
    ドライブ内で指定したフォルダID(またはルートID)の下に新たなフォルダを作成します。
  - `GetDriveRootItems(ctx context.Context, driveID string) ([]models.DriveItemable, error)`  
    指定ドライブのルートに含まれるアイテムを取得します。
  - `GetDriveItem(ctx context.Context, driveID, driveItemID string) (models.DriveItemable, error)`  
    指定したドライブアイテムの詳細情報を取得します。
  - `DeleteDriveItem(ctx context.Context, driveID, driveItemID string) error`  
    指定したアイテムを削除します。
  - `DownloadDriveItem(ctx context.Context, driveID, driveItemID string) ([]byte, error)`  
    ファイルの場合、コンテンツをダウンロードします。フォルダの場合はエラーが返されます。

- **サイト検索**
  - `GetSiteByName(ctx context.Context, siteName string) ([]models.Siteable, error)`  
    サイト名を検索して一致する SharePoint サイトを返します。

- **大容量ファイルアップロード**
  - `CreateUploadSession(ctx context.Context, driveID, itemPath string) (models.UploadSessionable, error)`  
    アップロードセッションを作成し、大容量ファイルの分割アップロードを可能にします。
  - `UploadFile(uploadSession models.UploadSessionable, file *os.File) (models.DriveItemable, error)`  
    取得したアップロードセッションを用いてファイルを分割アップロードします。

---

## ドメインレベルの処理 (service パッケージ)

`graphhelper` パッケージは、Microsoft Graph API への各種リクエストをシンプルなラッパー関数として提供しています。一方、ビジネスロジックやドメイン固有の検証は **service** パッケージで実装されています。  
この設計により、低レベルな API 呼び出しとドメインロジックの分離が実現され、利用者はシンプルなインターフェースを通じて各種機能を利用できるようになっています。

---

### サンプルコード (DownloadDriveItem の場合)

```go
func (s *graphServiceImpl) DownloadDriveItem(ctx context.Context, driveID, driveItemID string) ([]byte, error) {
    // まず対象のドライブアイテムを取得
    item, err := s.helper.GetDriveItem(ctx, driveID, driveItemID)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve drive item: %w", err)
    }
    // ファイルかどうかチェック（フォルダの場合は GetFile() が nil となる）
    if item.GetFile() == nil {
        return nil, fmt.Errorf("download is only supported for files, not folders")
    }
    // ファイルの場合のみダウンロード処理を実行
    return s.helper.DownloadDriveItem(ctx, driveID, driveItemID)
}
```
---

## 注意点

- Microsoft Graph の仕様変更や、依存ライブラリの更新により、インターフェースや挙動が変わる場合があります。
- 大容量ファイルアップロード時には、ネットワーク切断やセッションの有効期限切れに注意が必要です。本ライブラリでは再開機能を含む実装を行っていますが、必ずしもすべてのケースをカバーできる保証はありません。

---

## ライセンス

このライブラリは [MIT License](./LICENSE) で公開されています。詳細は LICENSE ファイルをご確認ください。

---

## 貢献・お問い合わせ

- バグ報告や新機能追加などの要望は、Issue や Pull Request で歓迎します。
- ご質問やご不明点がありましたら、GitHub の Issue を通じてお知らせください。
