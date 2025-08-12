# ビルド用のコンテナ
FROM golang:1.23-bullseye AS deploy-builer

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

#--------------------------------
# デプロイ用のコンテナ（本番環境でのアプリ起動）
FROM debian:bullseye-slim AS deploy

RUN apt-get update

COPY --from=deploy-builer /app/app .

EXPOSE 8080

CMD ["./app", "8080"]

#--------------------------------
# ローカル開発環境で利用するホットリロード環境
FROM golang:1.23-bullseye AS dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.40.4

CMD ["air"]