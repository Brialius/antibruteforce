FROM golang:1.13 AS builder

COPY . /app
WORKDIR /app
RUN make build-server

FROM alpine:3.7

COPY --from=builder /app/bin/server /app/server
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /app