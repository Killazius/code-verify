FROM golang:1.23

RUN apt-get update && apt-get install -y python3 python3-pip

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/web-server/
CMD ["./app"]
