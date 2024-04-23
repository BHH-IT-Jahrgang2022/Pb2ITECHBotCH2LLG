FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /logging .

EXPOSE 8080

CMD ["/logging"]
