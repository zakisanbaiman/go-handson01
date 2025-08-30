# 開発環境情報

## システム要件
- **OS**: Darwin (macOS)
- **Go**: 1.23.0以上
- **Docker**: Docker Compose v3.9対応

## ローカル開発セットアップ

### 1. 依存関係
- MySQL 8.0.29 (Dockerで提供)
- mysqldef (マイグレーションツール)

### 2. 環境変数
```bash
TODO_ENV=dev
PORT=8080
TODO_DB_HOST=todo-db
TODO_DB_PORT=3306
TODO_DB_USER=todo
TODO_DB_PASSWORD=todo
TODO_DB_NAME=todo
```

### 3. ポート設定
- アプリケーション: `localhost:18080`
- MySQL: `localhost:33306`

### 4. 起動手順
```bash
# 1. Dockerイメージビルド
make build-local

# 2. サービス起動
make up

# 3. マイグレーション実行
make migrate

# 4. 動作確認
curl http://localhost:18080/
```

## 開発ツール
- Air: ホットリロード対応 (設定: .air.toml)
- golangci-lint: 静的解析
- mysqldef: スキーママイグレーション

## ディレクトリマウント
- プロジェクトルート → `/app` (Docker内)
- MySQL設定 → `/etc/mysql/conf.d`