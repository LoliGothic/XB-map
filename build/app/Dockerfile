FROM golang:1.23-alpine

RUN apk upgrade --update && \
    apk --no-cache add git

WORKDIR /go/src/app
COPY . /go/src/app/

COPY go.mod /go/src/app/go.mod
COPY go.sum /go/src/app/go.sum

# modファイルにあるパッケージをダウンロード
RUN go mod download

# air（開発用）
RUN go install github.com/air-verse/air@v1.52.3

# airコマンドでgoを起動
CMD ["air", "-c", ".air.toml"]