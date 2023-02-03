# Build a goland docker image
FROM golang:1.19.5-alpine3.17
WORKDIR /app
COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o main .
ENTRYPOINT "/app/main"
