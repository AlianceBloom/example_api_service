FROM golang

ADD . /go/src/github.com/aliancebloom/example_api_service


RUN go get ./...
RUN go install github.com/aliancebloom/example_api_service

ENTRYPOINT /go/bin/example_api_service


EXPOSE 8080
