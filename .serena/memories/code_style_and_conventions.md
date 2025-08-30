# コードスタイル・規約

## Linting設定
- golangci-lint使用
- 有効なlinter:
  - goimports - インポート整理
  - unused - 未使用コード検出
  - errcheck - エラーハンドリングチェック
  - gocognit - 認知複雑度チェック
  - gocritic - コード品質チェック
  - gocyclo - 循環複雑度チェック (min-complexity: 30)
  - gofmt - フォーマット
  - govet - 静的解析
  - misspell - スペルチェック (US locale)
  - staticcheck - 静的解析
  - whitespace - 空白文字チェック

## 命名規則
- パッケージ名: 小文字、短縮形 (handler, entity, store, config)
- 構造体: パスカルケース (Server, Task)
- 関数・メソッド: パスカルケース (公開), キャメルケース (非公開)
- 定数: パスカルケース

## ディレクトリ構造規則
- handler/ - HTTPハンドラー
- entity/ - ドメインエンティティ
- store/ - データアクセス層
- config/ - 設定管理
- testutil/ - テスト共通処理
- testdata/ - テストデータ

## テスト規約
- `*_test.go` ファイル名
- テスト関数名: `TestXxx`
- テストオプション: `-race -shuffle=on`