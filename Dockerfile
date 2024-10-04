FROM golang:1.23.1

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd/app
RUN GOOS=linux go build -o /app/main

EXPOSE 8080
WORKDIR /app
CMD ["/app/main"]
