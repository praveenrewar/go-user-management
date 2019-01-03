FROM golang:1.11

WORKDIR /go/src/app
COPY . .

RUN export GO111MODULE=on
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["golang-mvc-boilerplate"]