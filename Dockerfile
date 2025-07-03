FROM alpine:3.18.4 AS certs

RUN apk add -U --no-cache ca-certificates && \
    addgroup -g 1001 app && \
    adduser app -u 1001 -D -G app /home/app

FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/main.go

FROM scratch AS final

COPY --from=certs /etc/passwd /etc/passwd
COPY --from=certs /etc/group /etc/group
COPY --chown=1001:1001 --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder /app/server /server

USER app

ENTRYPOINT ["/server"]