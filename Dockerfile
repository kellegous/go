# Use an intermediate container for initial building
FROM golang:1.14.0-buster as builder
RUN apt update && apt install -y upx-ucl && apt-get clean && rm -rf /var/lib/apt/lists/*

# Use go modules and let go packages call C code
ENV GO111MODULE=on CGO_ENABLED=1
WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-s -w" -o shorty ./cmd/shorty

# Move the built binary to /dist
WORKDIR /dist
RUN cp /build/shorty ./shorty
# this will collect dependent libraries so they're later copied to the final image
# NOTE: make sure you honor the license terms of the libraries you copy and distribute
# h/t @vaind/Ivan Dlugos: https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i
RUN ldd shorty | tr -s '[:blank:]' '\n' | grep '^/' | xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/
# Compress the binary and verify the output using UPX
# h/t @FiloSottile/Filippo Valsorda: https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
RUN upx --ultra-brute shorty && upx -t shorty
RUN mkdir /data

# Copy the contents of /dist to the root of a scratch containter
FROM scratch
COPY --chown=0:0 --from=builder /dist /
USER 65534
WORKDIR /
EXPOSE 8067
ENTRYPOINT ["/shorty"]
