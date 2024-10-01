
FROM golang:1.23.1 

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd/app
RUN GOOS=linux go build -o /app/main



ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USERNAME=${DB_USERNAME}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_DBNAME=${DB_DBNAME}
EXPOSE 8080
WORKDIR /app
CMD ["/app/main"]
