FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o main cmd/api/main.go

CMD ["./main"]
