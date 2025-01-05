FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o substantiate .

FROM alpine:3.21

COPY --from=builder /app/substantiate /substantiate

ENTRYPOINT ["/substantiate"]