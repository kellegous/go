FROM golang:alpine

COPY . /src
RUN cd /src/cmd/go && go build -mod=vendor -o /usr/bin/go

EXPOSE 8067

CMD ["/usr/bin/go", "--data=/data"]
