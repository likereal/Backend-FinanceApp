# Use Go 1.24.x Alpine-based image
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

EXPOSE 1323

CMD ["./server"]
