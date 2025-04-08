FROM golang:1.24-alpine

WORKDIR /app

COPY main.go .
COPY go.mod .

RUN go build -o main .

CMD ["./main"]
