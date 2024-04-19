FROM golang:latest 

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

EXPOSE 8080

ENV MATCHER_URL="http://host.docker.internal:88"

CMD ["/app/main"]