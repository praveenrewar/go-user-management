FROM golang:1.11.4-alpine

WORKDIR /go/src/golang-mvc-boilerplate
COPY . .

RUN apk update && apk upgrade && \
  apk add --no-cache bash git openssh && \
  export GO111MODULE=on && go mod vendor

RUN go build main.go

CMD ["./main"]