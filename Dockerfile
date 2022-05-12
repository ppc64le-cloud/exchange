FROM golang:1.18 as builder
WORKDIR /go/src/github.com/ppc64le-cloud/exchange
COPY cmd cmd
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum
COPY Makefile Makefile

RUN make all

FROM alpine:3.15

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/ppc64le-cloud/exchange/bin/ ./

CMD ["./pac-server"]
