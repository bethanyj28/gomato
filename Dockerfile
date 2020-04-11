FROM golang:latest

WORKDIR /app/

ADD . /go/src/github.com/bethanyj28/gomato

RUN go build -o gomato /go/src/github.com/bethanyj28/gomato/cmd/server/main.go
RUN cp /go/src/github.com/bethanyj28/gomato/environment.env .

ENTRYPOINT ["./gomato"]

EXPOSE 8080
