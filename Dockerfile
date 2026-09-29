FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/tradeshield-api ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/tradeshield-api /app/tradeshield-api

EXPOSE 8080
ENV PORT=8080

CMD ["/app/tradeshield-api"]
