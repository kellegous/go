FROM alpine

ENV GOPATH /go
COPY . /go/src/github.com/HALtheWise/o-links
RUN apk update \
  && apk add go git musl-dev \
  && go get github.com/HALtheWise/go-links \
  && apk del go git musl-dev \
  && rm -rf /var/cache/apk/* \
  && rm -rf /go/src /go/pkg \
  && mkdir /data

# Port might be dynamically set by Heroku, this may or may not be a problem
EXPOSE 8067

CMD ["/go/bin/o-links", "--data=/data"]
