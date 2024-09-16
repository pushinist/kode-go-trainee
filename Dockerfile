FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

ENV CGO_ENABLED=1

COPY . .

RUN go mod download

RUN go build -o /goapp

FROM alpine:latest

WORKDIR /app

COPY --from=builder /goapp .

COPY templates/ ./templates/

COPY static/ ./static/

EXPOSE 8080

CMD ["./goapp"]