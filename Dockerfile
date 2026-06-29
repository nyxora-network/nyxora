FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o nyxora ./cmd/nyxora

FROM alpine:3.19
RUN apk add --no-cache bash iproute2 wireguard-tools openssh-client sshpass
COPY --from=builder /build/nyxora /usr/local/bin/nyxora
RUN mkdir -p /etc/nyxora /var/log/nyxora
ENTRYPOINT ["nyxora"]
