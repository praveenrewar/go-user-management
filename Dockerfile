FROM golang:1.11

WORKDIR /go/src/golang-mvc-boilerplate
COPY . .

RUN export GO111MODULE=on && go mod vendor
RUN go build main.go

CMD ["./main"]