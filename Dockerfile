FROM golang:1.20-alpine3.16 AS build

RUN apk update && apk add gcc musl-dev gcompat libc-dev linux-headers
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bfs-portal *.go

FROM alpine:3.16
EXPOSE 8100

WORKDIR /app

COPY --from=build /app/bfs-portal /app/bfs-portal

RUN chmod +x /app/bfs-portal
ENTRYPOINT [ "/app/bfs-portal" ]
