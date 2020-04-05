FROM golang:alpine

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /rig

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN build.sh dev

EXPOSE 8080

CMD ["/rig/bin/upspeak"]
