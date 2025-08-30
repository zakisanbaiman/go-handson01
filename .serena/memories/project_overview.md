# Go Handson01 プロジェクト概要

## プロジェクトの目的
- GoでTODOアプリケーションを作成するハンズオンプロジェクト
- Web APIサーバーとMySQLデータベースを使用したタスク管理システム

## 技術スタック
- **言語**: Go 1.23.0
- **Webフレームワーク**: Chi router (github.com/go-chi/chi)
- **データベース**: MySQL 8.0.29
- **ORM**: sqlx (github.com/jmoiron/sqlx)
- **バリデーション**: go-playground/validator/v10
- **設定管理**: env/v6
- **コンテナ**: Docker & Docker Compose

## プロジェクト構造
```
.
├── main.go              # エントリーポイント
├── server.go            # サーバー構造体と起動処理
├── mux.go              # ルーティング設定
├── handler/            # HTTPハンドラー
│   ├── add_task.go
│   ├── list_task.go
│   └── handlers.go
├── entity/             # エンティティ定義
│   └── task.go
├── store/              # データストア層
│   ├── store.go
│   ├── repository.go
│   └── task.go
├── config/             # 設定管理
├── testutil/           # テストユーティリティ
├── clock/              # 時間関連ユーティリティ
└── _tools/mysql/       # データベーススキーマとMySQL設定
```

## アーキテクチャパターン
- レイヤードアーキテクチャ
  - Handler層: HTTPリクエスト処理
  - Entity層: ドメインオブジェクト
  - Store層: データアクセス層