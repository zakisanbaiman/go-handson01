# 推奨コマンド一覧

## 開発コマンド

### Docker関連
- `make build-local` - ローカル開発用Dockerイメージをビルド
- `make up` - Dockerコンテナを起動
- `make down` - Dockerコンテナを停止
- `make logs` - ログを表示
- `make ps` - コンテナ状況を表示

### テスト
- `make test` - テストを実行 (`go test -race -shuffle=on ./...`)
- `go test ./handler` - 特定パッケージのテスト実行

### データベース
- `make migrate` - マイグレーション実行
- `make dry-migrate` - マイグレーションのドライラン

### リンティング・フォーマット
- `golangci-lint run` - リンターを実行
- `gofmt -s -w .` - コードフォーマット
- `goimports -w .` - インポート整理

### 直接実行
- `go run .` - アプリケーション直接起動
- `go build` - バイナリビルド

### その他のユーティリティコマンド (Darwin)
- `ls -la` - ディレクトリ内容表示
- `find . -name "*.go"` - Goファイル検索
- `grep -r "pattern" .` - 文字列検索
- `git status` - Git状況確認