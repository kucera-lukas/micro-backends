FROM golang:1.18 as builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" -v -o postgres cmd/postgres/postgres.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/postgres ./postgres

ENTRYPOINT ["./postgres"]
