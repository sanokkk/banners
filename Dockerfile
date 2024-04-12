FROM golang:1.22-rc AS builder

COPY ./ /src/build

WORKDIR /src/build

COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN go build -o /src/build/exec /src/build/cmd/banners/main.go

FROM golang:1.22-rc AS production

COPY --from=builder src/build/config /config
COPY --from=builder /src/build/exec /build/banners
COPY --from=builder /src/build/.env /

WORKDIR /

EXPOSE 8080

CMD ["/build/exec"]