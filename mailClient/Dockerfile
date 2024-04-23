FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o mail .

EXPOSE 8080

CMD ["./mail"]
