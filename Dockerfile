FROM golang:1.22-rc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/banners ./cmd/banners/main.go

CMD ["./bin/banners"]