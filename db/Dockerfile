FROM golang

WORKDIR /app
COPY . . 

RUN go mod download

RUN go build -o /db

EXPOSE 8080

CMD ["/db"]