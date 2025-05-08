FROM golang:1.24.2-alpine AS builder

COPY . .

RUN go build -o /app ./...

FROM alpine:3.20

RUN apk add curl

COPY --from=builder /app /app

ENTRYPOINT ["/app"]