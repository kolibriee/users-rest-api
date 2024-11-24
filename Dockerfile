FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN GOOS=linux CGO_ENABLED=0 go build -o app -ldflags '-s -w' /app/cmd/app/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/app ./app
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/migrations ./migrations

CMD ["./app"]