FROM golang

WORKDIR /app
COPY . . 

RUN go mod download

RUN go build -o /matcher

EXPOSE 8081

CMD ["/matcher"]
