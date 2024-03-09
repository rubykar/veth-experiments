FROM golang:latest

WORKDIR /app

COPY client.go .

RUN go build -o main client.go

CMD ["./client"]
