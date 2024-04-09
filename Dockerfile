FROM golang:1.19 AS builder

RUN mkdir -p /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make PREFIX=/bin

FROM  quay.io/prometheus/busybox:latest
LABEL maintainer "Lukas Malkmus <mail@lukasmalkmus.com>"

COPY --from=builder /bin/faktory_exporter /bin/faktory_exporter

EXPOSE      9386
ENTRYPOINT  ["/bin/faktory_exporter"]
