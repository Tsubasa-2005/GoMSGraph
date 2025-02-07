# GoMSGraph ドキュメント

[![Build Status](https://github.com/Tsubasa-2005/GoMSGraph/actions/workflows/ci.yml/badge.svg)](https://github.com/Tsubasa-2005/GoMSGraph/actions)
[![Coverage Status](https://codecov.io/gh/Tsubasa-2005/GoMSGraph/branch/main/graph/badge.svg)](https://codecov.io/gh/Tsubasa-2005/GoMSGraph)

## 概要

このパッケージは、Microsoft Graph API との連携を抽象化するサービス層を提供します。  
内部では、公式 SDK を利用した `graphhelper` と、公式 SDK でサポートされていない機能（例：Chunk Upload 等）を実現するための独自実装 `httpclient` を組み合わせています。  
このサービス層は、`GraphService` インターフェースを介して利用でき、利用者は内部実装の詳細に依存することなく機能を利用可能です。

## カスタム HTTP クライアントと将来的な移行について

現時点では、HTTP クライアント (`httpclient`) のレスポンスは JSON を Go の独自型に変換して利用しています。  
たとえば、`GetDriveRootChildrenItemsRes` や `UploadSimpleFileRes` といった型をそのまま使用していますが、  
**将来的に公式の Microsoft Graph SDK が同様の機能（Chunk Upload など）をサポートし、公式の `models` に対応する型が提供された際には、これらの独自型を公式型へ置き換える予定です。**

この置き換えをスムーズに実施できるよう、サービス層のインターフェース (`GraphService`) を通して利用者に機能を提供しており、  
内部のレスポンス型への直接依存を避ける設計としています。

## ディレクトリ構成

- **graphhelper/**
    - 公式 SDK（`msgraph-sdk-go`）と Azure Identity を利用して、Graph API の認証やサイト検索などを行うラッパを実装しています。

- **httpclient/**
    - 公式 SDK で未対応の機能（Chunk Upload 等）を実現するための独自実装。
    - HTTP レスポンスの JSON をパースして、独自に定義した型（例：`GetDriveRootChildrenItemsRes`、`UploadSimpleFileRes` など）を返します。

- **service/**
    - `GraphService` インターフェースを定義し、上記コンポーネントを組み合わせたビジネスロジックを実装しています。

## 利用例

以下は、`GraphService` を利用した簡単な利用例です:

```go
package main

import (
  "context"
  "fmt"
  "os"

  "github.com/Tsubasa-2005/GoMSGraph/service"
)

func main() {
  clientId := os.Getenv("CLIENT_ID")
  tenantId := os.Getenv("TENANT_ID")
  clientSecret := os.Getenv("CLIENT_SECRET")
  baseURL := "https://graph.microsoft.com/v1.0"

  // GraphService のインスタンス作成
  gs, err := service.NewGraphService(clientId, tenantId, clientSecret, baseURL)
  if err != nil {
    panic(err)
  }

  ctx := context.Background()

  // サイトを名前で検索
  sites, err := gs.GetSiteByName(ctx, "contoso")
  if err != nil {
    panic(err)
  }
  fmt.Println(sites)
}
```

## 今後の計画

- **公式 SDK への移行:**  
  公式 SDK が本プロジェクトで実現している機能（例：ファイルアップロード、Chunk Upload 等）をサポートし、公式の `models` が利用可能になった場合、  
  現在の `httpclient` 内で定義している独自型を公式の型へ切り替えます。

- **リファクタリング:**  
  利用者が `GraphService` インターフェースを通じて機能を利用できるよう、内部実装の変更が外部に影響しないよう努めます。

- **エラーハンドリングおよびページング対応の改善:**  
  HTTP ステータスコードのより柔軟な扱いや、大量データに対応するためのページング機能の実装を検討しています。
