FROM golang:alpine
MAINTAINER Amin Jams <aminjam.software@gmail.com>

RUN apk update && apk add git bash

RUN go get github.com/aminjam/goflat/cmd/goflat
RUN echo "$(goflat --version)"

CMD ["bash"]
