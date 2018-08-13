FROM golang:1.10

RUN apt-get update
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt-get install -y nodejs closure-compiler
RUN npm install -g typescript sass

WORKDIR /go/src/github.com/kellegous/go
COPY . .
RUN make ALL
RUN go get -u -d github.com/kellegous/go
RUN CGO_ENABLED=0 go build -v -o go .

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/kellegous/go/go .
CMD ["./go"]
