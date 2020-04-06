FROM golang:1.14

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

WORKDIR /go/src/upspeak

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o bin/upspeak .

FROM alpine:latest
WORKDIR /
COPY --from=0 /go/src/upspeak/bin/upspeak .
ENV ENV=dev
EXPOSE 8080
CMD ["/upspeak"]
