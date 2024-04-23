FROM golang

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o /unsolveddb

EXPOSE 8080

CMD ["/unsolveddb"]