FROM golang:latest

WORKDIR /app

COPY /logging/ .
COPY main.go .

RUN go mod download

CMD ["go", "run", "main.go"]