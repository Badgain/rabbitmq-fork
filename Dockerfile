FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod tidy

RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

RUN goimports -w .
RUN golangci-lint run

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app

# Final stage
FROM scratch
WORKDIR /app
# Copying the built app and configurations from builder stage
COPY --from=builder /app /app
COPY --from=builder /mappings /mappings
ENTRYPOINT ["/app"]