FROM golang:latest

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

WORKDIR /app

ENV SRC_DIR=/go/src/github.com/bethanyj28/gomato
ADD . $SRC_DIR
RUN cd $SRC_DIR; go build -o gomato cmd/server/main.go; cp gomato /app/; cp environment.env /app/

ENTRYPOINT ["./gomato"]
