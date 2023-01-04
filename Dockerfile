FROM golang:1.17-alpine

COPY . /app

WORKDIR /app
RUN go build -o /app/udpService

CMD ["/app/udpService"]