FROM golang:1.11.4-alpine

WORKDIR /go/src/go-user-management
COPY . .

RUN apk update && apk upgrade && \
  apk add --no-cache bash git openssh build-base && \
  export GO111MODULE=on && go mod init && go mod vendor

RUN go test -c -coverpkg ./...
RUN go build .

CMD echo "Running Unit Tests" && WorkEnv=test "./go-user-management.test" && WorkEnv=dev "./go-user-management"
