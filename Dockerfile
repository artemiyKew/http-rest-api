FROM golang:1.21.3 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Builder
FROM golang:1.21.3 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/apiserver

# Final
FROM scratch
COPY --from=builder /app/configs /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]