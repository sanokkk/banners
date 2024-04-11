FROM golang:1.22-rc AS builder

COPY ./ /src/build

WORKDIR /src/build

COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN go build -o /src/build/exec /src/build/cmd/banners/main.go

FROM golang:1.22-rc AS production

COPY --from=builder src/build/config/docker.yaml /config
COPY --from=builder *exe /build/banners

WORKDIR /

EXPOSE 8080

CMD ["exec"]