FROM golang:1.22.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/myapp ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add tzdata ca-certificates
ENV TZ=America/Santiago

RUN addgroup -S app -g 1000 && adduser -S -g app app --uid 1000

COPY --from=builder --chown=app:app /app/myapp /app/myapp

USER app
WORKDIR /app

EXPOSE 8080
ENTRYPOINT ["/app/myapp"]
