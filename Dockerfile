FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gostash

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gostash .
COPY . .

RUN chmod +x ./gostash

EXPOSE 8000/tcp

CMD ["./gostash"]

