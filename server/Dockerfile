FROM golang:1.23
LABEL authors="Broderick-Westrope"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./server ./server
COPY ./internal ./internal

WORKDIR /app/server

RUN go build -o server .

EXPOSE 8080

ENV SERVER_ADDR=$SERVER_ADDR
ENV REDIS_ADDR=$REDIS_ADDR
ENV DB_CONN=$DB_CONN
ENV LOG_LEVEL=$LOG_LEVEL

CMD ["./server"]
