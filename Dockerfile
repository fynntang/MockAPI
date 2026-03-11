FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mockapi ./cmd/mockapi

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/mockapi /usr/local/bin/mockapi
EXPOSE 8088
ENTRYPOINT ["mockapi"]
CMD ["8088"]
