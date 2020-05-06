FROM golang:1.11 AS builder

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p $GOPATH/src/github.com/lukasmalkmus/faktory_exporter/

COPY . $GOPATH/src/github.com/lukasmalkmus/faktory_exporter/

WORKDIR $GOPATH/src/github.com/lukasmalkmus/faktory_exporter/

RUN dep ensure

RUN make PREFIX=/bin

FROM  quay.io/prometheus/busybox:latest
LABEL maintainer "Lukas Malkmus <mail@lukasmalkmus.com>"

COPY --from=builder /bin/faktory_exporter /bin/faktory_exporter

EXPOSE      9386
ENTRYPOINT  ["/bin/faktory_exporter"]