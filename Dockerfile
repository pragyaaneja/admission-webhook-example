FROM golang:1.15-buster as builder

ADD . /build
WORKDIR /build

RUN go get ./...
RUN make

FROM alpine

COPY --from=builder /build/initc   /usr/local/bin/initc
COPY --from=builder /build/webhook /usr/local/bin/webhook
