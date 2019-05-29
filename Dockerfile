FROM golang:1.11.4-alpine

WORKDIR /go/src/datax-dp
COPY . .

RUN apk update && apk upgrade && \
  apk add --no-cache bash git openssh build-base && \
  export GO111MODULE=on && go mod vendor

RUN go test -c -coverpkg ./...
RUN go build .

CMD echo "Running Unit Tests" && WorkEnv=test "./datax-dp.test" && WorkEnv=dev "./datax-dp"
