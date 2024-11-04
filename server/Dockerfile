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

# ENV REDIS_ADDRESS=redis:6379

CMD ["./server"]
