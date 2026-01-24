FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./main.go

FROM alpine:latest

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/internal/pkg/email/template /app/internal/pkg/email/template

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

USER appuser

EXPOSE 8080

CMD ["/app/main"]
