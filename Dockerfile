FROM golang:1.20.5-buster as builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/go-proxypool ./cmd/main.go

FROM debian:buster-slim

RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /bin/go-proxypool /bin/go-proxypool
COPY config.yaml.sample /etc/proxypool/config.yaml
CMD ["/bin/go-proxypool"]