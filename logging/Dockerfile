FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /logger .

EXPOSE 8080

CMD ["/logger"]
