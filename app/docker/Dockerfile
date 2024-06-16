FROM golang:1.17-alpine AS gobuilder
WORKDIR /go/src/store-api


COPY go.mod .
COPY go.sum .
COPY .env .

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

COPY . .

ENTRYPOINT CompileDaemon --build="go build -o /go/bin/paylater-customer-api ." --command=/go/bin/paylater-customer-api --directory=.

