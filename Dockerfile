FROM alpine

ENV GOPATH /go
COPY . /go/src/github.com/HALtheWise/go-links
RUN apk update \
  && apk add go git musl-dev \
  && go get github.com/HALtheWise/go-links \
  && apk del go git musl-dev \
  && rm -rf /var/cache/apk/* \
  && rm -rf /go/src /go/pkg \
  && mkdir /data

EXPOSE 8067

CMD ["/go/bin/go-links", "--data=/data"]
